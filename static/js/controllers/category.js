angular.module('categoryApp', []).controller('categoryCtrl', function($scope, $http) {
	$http.get("/category/load").success(function(data) {
		$scope.Categories = data;
	});

});
