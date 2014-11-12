var app = angular.module('app', ['ngResource']);

app.controller('MainController', function($scope, $resource) {
    var api = $resource('api/expenses');

    // Hardcoded settings for now
    $scope.users = [ "David", "Sofie" ];
    // TODO cateogires with default amounts and allowance per user
    $scope.categories = [ "Mat", "Gemensamt", "Eget inkop", "Betalning" ];

    $scope.expenses = api.query({user: 1});

    $scope.newExpense = {
        user: 1, // TODO save last selected user
        category: $scope.categories[0], // TODO save last selected
        date: "2014-01-01T01:01:01Z" 
    }

    $scope.submitExpense = function() {
        api.save($scope.newExpense, {});
    }
});
