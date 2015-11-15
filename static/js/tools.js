/* * 用下面的代码就不会发生悲剧了
* 最终结论
* E(e).stop(); 阻止事件冒泡
* E(e).prevent();阻止事件默认行为
*/
var E = function(e){
	e = window.event || e;
	return {
			stop: function() {
				if (e && e.stopPropagation) e.stopPropagation();
				else e.cancelBubble = true
			},
			prevent: function() {
				if (e && e.preventDefault) e.preventDefault();
				else e.returnValue = false
			}
	}
}