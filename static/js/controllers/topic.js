angular.module('topicApp', []).controller('topicCtrl', function($scope, $http) {
	$http.get("/topic/load").success(function(data) {
		$scope.Topics = data;
	});
	 
});


angular.module('topic_modifyApp', []).controller('topic_modifyCtrl', function($scope, $http) {
	var tid = $("#tid").val();
	$http.get("/topic/loadModify?tid=" + tid).success(function(data) {
		$scope.Topic = data;
	});
	 
});



angular.module('topic_viewApp', []).controller('topic_viewCtrl', function($scope, $http) {
	var tid = $("#tid").val();
	$http.get("/topic/loadView?tid=" + tid).success(function(data) {
		$scope.Topic = data.Topic;
		$scope.Lables = data.Lables;
		$scope.Replies = data.Replies;
	});
	 
});
