$(document).ready(function() {
	$('.masthead').visibility({
		once: false
	});
	$('.image').visibility({
		type: 'image',
		transition: 'vertical flip in',
		duration: 500
	});
	$('.masthead').removeClass('zoomed');
	$('.main.menu  .ui.dropdown').dropdown({
		on: 'hover'
	});
	if ($(window).width() > 600) {
		return $('body').visibility({
			offset: -10,
			observeChanges: false,
			once: false,
			continuous: false,
			onTopPassed: function() {
				return requestAnimationFrame(function() {
					$('.following.bar').addClass('light fixed').find('.menu').removeClass('inverted');
					return $('.following .additional.item').transition('scale in', 350);
				});
			},
			onTopPassedReverse: function() {
				return requestAnimationFrame(function() {
					return $('.following.bar').removeClass('light fixed').find('.menu').addClass('inverted').find('.additional.item').transition('hide');
				});
			}
		});
	}
});