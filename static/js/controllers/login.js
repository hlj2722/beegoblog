angular.module('loginApp', []).controller('loginCtrl', function($scope, $http) {
	$scope.backToHome = function () {
      		window.location.href = "/";
     }
	
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


