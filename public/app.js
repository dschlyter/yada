var app = angular.module('app', ['ngResource']);

app.controller('MainController', function($scope, $resource) {
  var init = function() {
    $scope.expenses = [];

    // Hardcoded settings for now
    $scope.users = [ 
      {id: 1, name: "David"}, 
      {id: 2, name: "Sofie"}
    ];

    $scope.categories = [ 
      { title: "Mat", users: [1, 2], split: [60, 40] }, 
      { title: "Gemensamt", users: [1, 2], split: [50, 50] }, 
      { title: "Betalning", users: [1], split: [0, 100] }, 
      { title: "Betalning", users: [2], split: [100, 0] }, 
      { title: "Eget", users: [1], split: [100, 0] }, 
      { title: "Eget", users: [2], split: [0, 100] }
    ];

    $scope.availableCategories = [];

    $scope.user = loadDefault("userid", 1);
    $scope.category = $scope.categories[loadDefault("categoryIndex", 0)];
    $scope.submitDisabled = false;

    $scope.fetchLatestExpenses();
    initExpense();
    $scope.userChanged();
  }

  // TODO extract to data service?
  var api = $resource('api/expenses');
  var requestInProgress = false;
  var everythingFetched = false;

  $scope.fetchLatestExpenses = function() {
    fetchExpenses(undefined);
  }

  $scope.fetchMoreExpenses = function() {
    if (requestInProgress || everythingFetched) {
      return;
    }

    var afterKey = undefined;
    if ($scope.expenses.length > 0) {
      afterKey = $scope.expenses[$scope.expenses.length -1].Id;
    }
    fetchExpenses(afterKey);
  }

  var fetchExpenses = function(afterKey) {
    var query = {user: 1, limit: 20};

    if (afterKey !== undefined) {
      query.afterKey = afterKey;
    }
      
    requestInProgress = true;
    var newData = api.query(query, function() {
      var mergedList = [];

      var oldIndex = 0;
      var newIndex = 0;
      do {
        var oldId = oldIndex < $scope.expenses.length ? $scope.expenses[oldIndex].Id : "";
        var newId = newIndex < newData.length ? newData[newIndex].Id : "";

        // Merge largest id into list, skip new if ids are equal
        if (oldId > newId) {
          mergedList.push($scope.expenses[oldIndex]);
          oldIndex++;
        } else if (newId > oldId) {
          mergedList.push(newData[newIndex]);
          newIndex++;
        } else {
          newIndex++
        }
      } while (oldId !== "" || newId !== "");

      $scope.expenses = mergedList;
      requestInProgress = false;
      if (newData.length === 0) {
        everythingFetched = true;
      }
    }, function(error) {
      setError(error);
      requestInProgress = false;
    });
  }

  var setError = function(error) {
    $scope.error = "Connection Error"
    if (error && error.data && error.data.Error) {
      $scope.error = error.data.Error
    }
  }

  var initExpense = function() {
    $scope.newExpense = {
      user: $scope.user,
      category: $scope.category.title,
    }

    today = moment().startOf("day")
    $scope.browserDate = today.format("YYYY-MM-DD");
    $scope.mobileDate = today.toDate();
    $scope.dateChanged(today);

    if ($scope.form) {
        $scope.form.$setPristine()
        // TODO set focus on Total field - but no DOM-manipulation in controller, hmm
    }
  }

  $scope.dismissError = function() {
    $scope.error = null;
  }

  $scope.selectUser = function(user) {
    $scope.user = user;
    $scope.userChanged();
  }

  $scope.userChanged = function() {
    $scope.newExpense.user = $scope.user
    $scope.availableCategories = $scope.categories.filter(function(elem) {
      return elem.users.indexOf($scope.newExpense.user) > -1;
    });

    // User change may cause current category to be unavailable
    if ($scope.availableCategories.indexOf($scope.category) == -1) {
      // Switch to "mirror" category if possible 
      var mirror = $scope.availableCategories.filter(function(elem) {
        return elem.title === $scope.category.title;
      });
      $scope.category = mirror ? mirror[0] : $scope.availableCategories[0]; 
    }

    $scope.categoryChanged();

    saveDefault("userid", $scope.user);
  }

  $scope.selectCategory = function(category) {
    $scope.category = category;
    $scope.categoryChanged();
  }

  $scope.categoryChanged = function() {
    $scope.newExpense.category = $scope.category.title;
    var otherUser = $scope.newExpense.user == 1 ? 2 : 1;
    $scope.percentage = $scope.category.split[otherUser - 1];
    $scope.percentageChanged();

    var index = $scope.categories.indexOf($scope.category);
    if (index > -1) {
      saveDefault("categoryIndex", index);
    }
  }

  $scope.dateChanged = function(newDate) {
    $scope.newExpense.date = moment(newDate).format("YYYY-MM-DDTHH:mm:ss") + "Z";
  }

  // Four inputs depend on each other, recalc all when they change
  // Fixed values for calc are total and percentage
  // TODO unit tests for this
  var calcOwedFromPercentage = function() {
    var exp = $scope.newExpense;
    exp.owedAmount = Math.round(exp.totalAmount * ($scope.percentage / 100));
  }

  var calcPercentageFromOwed = function() {
    var exp = $scope.newExpense;
    $scope.percentage = Math.round(100 * exp.owedAmount / exp.totalAmount);
  }

  var calcOwedFromDiff = function() {
    var exp = $scope.newExpense;
    exp.owedAmount = Math.round((exp.totalAmount + $scope.diff) / 2);
  }

  var calcDiffFromOwed = function() {
    var exp = $scope.newExpense;
    $scope.diff = exp.owedAmount - (exp.totalAmount - exp.owedAmount);
  }

  $scope.totalChanged = function() {
    calcOwedFromPercentage();
    calcDiffFromOwed();
  }

  $scope.percentageChanged = function() {
    calcOwedFromPercentage();
    calcDiffFromOwed();
  }

  $scope.owedChanged = function() {
    calcPercentageFromOwed();
    calcDiffFromOwed();
  }

  $scope.diffChanged = function() {
    calcOwedFromDiff();
    calcPercentageFromOwed();
  }

  $scope.submitExpense = function() {
    $scope.dismissError();
    $scope.submitDisabled = true; // Avoid multiple submits before response
    api.save($scope.newExpense, {}, function() {
      initExpense();
      $scope.fetchLatestExpenses();
      $scope.submitDisabled = false;
    }, function(error) {
      setError(error);
      $scope.submitDisabled = false;
    });
  }

  $scope.isFuture = function(expense) {
    return expense && moment(expense.Date).isAfter(moment())
  }

  $scope.selectRow = function(row) {
    if ($scope.selectedRow === row) {
      $scope.selectedRow = null;
    } else {
      $scope.selectedRow = row;
    }
  }

  // TODO extract to service
  var loadDefault = function(key, defaultValue) {
    try {
      var ret = localStorage[key];
      if (ret !== undefined) {
        return parseInt(ret);
      }
    } catch(e) {
      // No op
    }

    return defaultValue;
  }

  var saveDefault = function(key, value) {
    try {
      localStorage[key] = value;
    } catch(e) {
      // No op
    }
  }

  // Call init last when everything has initialized
  init();
});
