
var Storage = function() {
	this.url = "/api/frameworks";

	this.load = function(page, cb) {
	  $.ajax({
		url: this.url + "?perPage=5&page=" + page,
		dataType: 'json',
		cache: false,
		success: function(data) {
			cb(data)
		}.bind(this),
		error: function(xhr, status, err) {
		  console.error(this.url, status, err.toString());
		}.bind(this)
	  });
	}

	this.add = function(framework, cb) {
		$.ajax({
		  url: this.url,
		  dataType: 'json',
		  type: 'POST',
		  data: JSON.stringify(framework),
		  success: function(data) {
			cb(data)
		  }.bind(this),
		  error: function(xhr, status, err) {
			console.error(this.url, status, err.toString());
		  }.bind(this)
		});
	}

	this.remove = function(frameworkID, cb) {
		$.ajax({
		  url: this.url + "/" + frameworkID,
		  dataType: 'json',
		  type: 'DELETE',
		  success: function(data) {
				cb(data)
		  }.bind(this),
		  error: function(xhr, status, err) {
				console.error(this.url, status, err.toString());
		  }.bind(this)
		});
	}
}
