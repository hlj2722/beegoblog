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
	
	$scope.checkInput = function ($event) {
    		var uname = $("#uname").val();
    		if (uname.length == 0) {
    			alert("请输入帐号");
            	E($event).stop();
			E($event).prevent();
    		}

    		var pwd = $("#pwd").val();
    		if (pwd.length == 0) {
    			alert("请输入密码");
            	E($event).stop();
			E($event).prevent();
    		}
    	}

	 
});


angular.module('user_addApp', []).controller('user_addCtrl', function($scope, $http) {

	$scope.checkInput = function ($event) {
    		var uname = $("#uname").val();
    		if (uname.length == 0) {
    			alert("请输入帐号");
            	E($event).stop();
			E($event).prevent();
    		}

    		var pwd = $("#pwd").val();
    		if (pwd.length == 0) {
    			alert("请输入密码");
            	E($event).stop();
			E($event).prevent();
    		}
    	}
});


angular.module('user_viewApp', []).controller('user_viewCtrl', function($scope, $http) {
	var uname = $("#uname").val();
	$http.get("/user/loadView?uname=" + uname).success(function(data) {
		$scope.User = data.User;
	});
	 
});
