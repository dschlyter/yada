var app = angular.module('app', ['ngResource']);

app.controller('MainController', function($scope, $resource) {
  var api = $resource('api/expenses');

  var init = function() {
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

    $scope.user = 1; // TODO save last selected user
    $scope.category = $scope.categories[0]; // TODO save last selected
    $scope.submitDisabled = false;

    refresh();
    initExpense();
    $scope.userChanged();

    // TODO should be in a directive
    new Pikaday({
      field: document.getElementById('pikaday-picker'),
      firstDay: 1,
      yearRange: [2000,2020]
    });
  }

  var refresh = function() {
    $scope.expenses = api.query({user: 1});
  }

  var initExpense = function() {
    $scope.newExpense = {
      user: $scope.user,
      category: $scope.category.title,
    }

    $scope.browserDate = moment().format("YYYY-MM-DD");
    $scope.mobileDate = new Date();
    $scope.dateChanged(new Date());

    if ($scope.form) {
        $scope.form.$setPristine()
        // TODO set focus on Total field - but no DOM-manipulation in controller, hmm
    }
  }

  $scope.dismissError = function() {
    $scope.error = null;
  }

  $scope.userChanged = function() {
    $scope.newExpense.user = $scope.user
    $scope.availableCategories = $scope.categories.filter(function(elem) {
      return elem.users.indexOf($scope.newExpense.user) > -1;
    });
    if ($scope.category) {
      $scope.categoryChanged();
    }
  }

  $scope.categoryChanged = function() {
    $scope.newExpense.category = $scope.category.title;
    var otherUser = $scope.newExpense.user == 1 ? 2 : 1;
    $scope.percentage = $scope.category.split[otherUser - 1];
    $scope.percentageChanged();
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
      refresh();
      $scope.submitDisabled = false;
    }, function(error) {
      $scope.error = error.data.Error;
      $scope.submitDisabled = false;
    });
  }

  $scope.selectRow = function(row) {
    if ($scope.selectedRow === row) {
      $scope.selectedRow = null;
    } else {
      $scope.selectedRow = row;
    }
  }

  // Call init last when everything has initialized
  init();
});
