angular.module('categoryApp', []).controller('categoryCtrl', function($scope, $http) {
	$http.get("/category/load").success(function(data) {
		$scope.Categories = data;
	});
	

	$scope.checkInput =	function ($event) {
    		var uname = $("#name").val();
    		if (uname.length == 0) {
    			alert("请输入分类名称");
            	E($event).stop();
			E($event).prevent();    			
    		}
    }

});
