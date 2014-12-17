// TODO Support a better infinite-scroll http://www.zdnet.com/article/google-tries-to-save-the-web-from-the-curse-of-infinite-scrolling/
app.directive('infiniteScroll', function($interval) {
  return {
    link: function(scope, element, attrs) {
      callbackName = attrs.infiniteScroll;

      var check = function() {
        if (document.body.scrollTop + window.innerHeight + 20 > document.body.scrollHeight) {
          scope[callbackName]();
        }
      }

      window.onscroll = function() {
        check();
      };

      $interval(check, 1000);
    }
  }
});
