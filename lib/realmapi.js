(function($, window) {

    var BASEURL = 'http://aloof.io:5353/';

    var _cnt = 0;
    function queue_request(obj) {
	var oc = obj.complete;
	
	obj.complete = function() {
	    if (oc) oc.apply(this, arguments);	    
	    _cnt = $(document).queue('ajax').length;
	    update_counter();
	    $(document).dequeue('ajax');
	}
	if (_cnt) {
	    $(document).queue('ajax', function(){ $.ajax(obj) });
	} else {
	    $.ajax(obj);
	}
	_cnt++;
	update_counter();
    }

    function update_counter() {
	$('#counter').text(_cnt).parent().toggle(!!_cnt);
    }


    function realmAPI(path, opts, extraopts, callback) {
	opts.ignore = Math.floor(1e3 + 9e3 * Math.random())
	var url = BASEURL + path + '?' + $.param(opts)

	if (typeof extraopts == 'function') {
	    callback = extraopts
	    extraopts = {}
	}

	if (extraopts.iframe) {
	    var ifr = document.createElement('iframe')
	    ifr.style.display = 'none'
	    document.body.appendChild(ifr)
	    $(ifr).load(function() {
		$(this).remove()
		if (typeof callback == 'function') callback()
	    })
	    ifr.src = url
	    return
	}

	queue_request({
	    url: url,
	    complete: callback
	})
    }

    window.realmAPI = realmAPI

})($, window)

