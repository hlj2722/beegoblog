angular.module('homeApp', []).controller('homeCtrl', function($scope, $http) {
	var category = $("#category").val()
	var lable = $("#lable").val()
	
	$http.get("/home/load?category=" + category + "&lable=" +  lable).success(function(data) {
		$scope.Topics = data.Topics;
		$scope.Categories = data.Categories;
	});
	 
});
