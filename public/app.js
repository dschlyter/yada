var app = angular.module('app', ['ngResource']);

app.controller('MainController', function($scope, $resource) {
  var api = $resource('api/expenses');

  // Hardcoded settings for now
  $scope.users = [ "David", "Sofie" ];
  // TODO cateogires with default amounts and allowance per user
  $scope.categories = [ "Mat", "Gemensamt", "Eget inkop", "Betalning" ];

  var refresh = function() {
    $scope.expenses = api.query({user: 1});
  }
  refresh();

  $scope.newExpense = {
    user: 1, // TODO save last selected user
    category: $scope.categories[0], // TODO save last selected
    date: "2014-11-25T01:01:01Z" 
  }

  $scope.dismissError = function() {
    $scope.error = null;
  }

  $scope.categoryChanged = function() {

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
    // TODO 0 is allowed + better form validation
    if (!$scope.newExpense.totalAmount || !$scope.newExpense.owedAmount) {
      $scope.error = "Total and owed must be defined";
      return;
    }

    $scope.dismissError()
      api.save($scope.newExpense, {}, function() {
        refresh();
      }, function(error) {
        $scope.error = error.data.Error;
      });
  }
});
