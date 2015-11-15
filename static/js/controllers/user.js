angular.module('userApp', []).controller('userCtrl', function($scope, $http) {
	$http.get("/user/load").success(function(data) {
		$scope.Users = data;
	});
	 
});


angular.module('user_modifyApp', []).controller('user_modifyCtrl', function($scope, $http) {
	var uname = $("#uname").val();
	$http.get("/user/loadModify?uname=" + uname).success(function(data) {
		$scope.User = data;
	});
	 
});



angular.module('user_viewApp', []).controller('user_viewCtrl', function($scope, $http) {
	var uname = $("#uname").val();
	$http.get("/user/loadView?uname=" + uname).success(function(data) {
		$scope.User = data.User;
	});
	 
});
