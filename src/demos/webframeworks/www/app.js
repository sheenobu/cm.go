
// FrameworkRow defines our javascript library information block.
var FrameworkRow = React.createClass({
	render: function() {
		desc = this.props.data.Description;
		if (desc.length > 100) {
			desc = desc.substr(0,100) + "..."
		}
		var onClick = this.props.onFrameworkDelete.bind(this, this.props.data.Id);
		var classNames = "btn-floating btn-small waves-effect secondary-content waves-light red";
		return (
			<li className="collection-item avatar">
				<i className="material-icons circle">folder</i>
				<span className="title"><b>{this.props.data.Name}</b></span>
				<p>{desc}</p>
				<a href={this.props.data.Url}>{this.props.data.Url}</a>
				<a onClick={onClick}  className={classNames}>	<i className="material-icons">delete</i></a>
			</li>
			)
	}
});

// AddForm implements the form to add a new javascript library to the system.
var AddForm = React.createClass({
	render: function() {
		return (
			<div id="addFormModal" className="modal">
				<div className="modal-content">
					<div className="row">
						<form className="col s12">
							<div className="row">
								<div className="input-field col s12">
									<input ref="name" id="name" type="text" className="validate"/>
									<label for="name">Name</label>
								</div>
							</div>
							<div className="row">
								<div className="input-field col s12">
									<input ref="description" id="description" type="text" className="validate"/>
									<label for="description">Description</label>
								</div>
							</div>
							<div className="row">
								<div className="input-field col s12">
									<input ref="url" id="url" type="text" className="validate"/>
									<label for="url">Project URL</label>
								</div>
							</div>
						</form>
					</div>
				</div>
				<div className="modal-footer">
					<a href="#!" onClick={this.addFramework} className="modal-action modal-close waves-effect waves-green btn">Add Framework</a>
					<a href="#!" onClick={this.cancel} className="left fmodal-action modal-close waves-effect waves-green btn">Cancel</a>
				</div>
			</div>
		);
	},

	// addFramework calls the onFrameworkSubmit callback and then resets the form
	addFramework: function() {

		var name = React.findDOMNode(this.refs.name).value.trim();
    var desc = React.findDOMNode(this.refs.description).value.trim();
		var url = React.findDOMNode(this.refs.url).value.trim();

	  if (!name || !desc || !url) {
			return;
		}

	  this.props.onFrameworkSubmit({name: name, description: desc, url: url});

	  React.findDOMNode(this.refs.name).value = '';
		React.findDOMNode(this.refs.description).value = '';
		React.findDOMNode(this.refs.url).value = '';

	 	return;
	},

	// cancel resets the form
	cancel: function() {
	  React.findDOMNode(this.refs.name).value = '';
		React.findDOMNode(this.refs.description).value = '';
		React.findDOMNode(this.refs.url).value = '';
		return
	},
})

var FrameworkTable = React.createClass({
	render: function() {
		var self = this;
		var nodes = this.props.data.Frameworks.map(function(row) {
			return (
				<FrameworkRow key={row.Id} data={row} onFrameworkDelete={self.props.onFrameworkDelete} />
			);
		});
		return (
			<ul className="collection">
				{nodes}
			</ul>
		);
	}
});

var Frameworks = React.createClass({
  loadLibrariesFromServer: function() {
    $.ajax({
      url: this.props.url + "?perPage=5&page=" + this.state.data.CurrentPage,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data})
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
   },
	getInitialState: function() {
		return {data: {
			Frameworks: [],
			CurrentPage: 0,
			PageCount: 0,
		}}
	},
	componentDidMount: function() {
		this.loadLibrariesFromServer();
	},
	handleAddFrameworkSubmit: function(framework) {
		$.ajax({
		  url: this.props.url,
		  dataType: 'json',
		  type: 'POST',
		  data: JSON.stringify(framework),
		  success: function(data) {
				this.loadLibrariesFromServer();
		  }.bind(this),
		  error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
		  }.bind(this)
		});
	},
	handleFrameworkDelete: function(frameworkId) {
		$.ajax({
		  url: this.props.url + "/" + frameworkId,
		  dataType: 'json',
		  type: 'DELETE',
		  success: function(data) {
				this.loadLibrariesFromServer();
		  }.bind(this),
		  error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
		  }.bind(this)
		});
	},
	handlePageChange: function(i) {
		var d = this.state.data;
		d.CurrentPage = i;
		this.setState(d);
		this.loadLibrariesFromServer();
	},
	render: function() {

		var nodes = []

		for(var i = 0; i != this.state.data.PageCount; i++) {
			if (i == this.state.data.CurrentPage) {
				nodes[i] = (<li key={i} onClick={this.handlePageChange.bind(this, i)} className="active"><a href="#!">{i+1}</a></li>)
			} else {
				nodes[i] = (<li key={i} onClick={this.handlePageChange.bind(this, i)} className="waves-effect"><a href="#!">{i+1}</a></li>)
			}
		}

		var prevIcon = []

		if(this.state.data.CurrentPage == 0) {
			prevIcon[0] = (<li className="disabled"><a href="#prev"><i className="material-icons">chevron_left</i></a></li>)
		}else{
			prevIcon[0] = (<li className="waves-effect" onClick={this.handlePageChange.bind(this, this.state.data.CurrentPage-1)}><a href="#prev"><i className="material-icons">chevron_left</i></a></li>)
		}

		var nextIcon = []

		if(this.state.data.CurrentPage == this.state.data.PageCount-1) {
			nextIcon[0] = (<li className="disabled"><a href="#next"><i className="material-icons">chevron_right</i></a></li>)
		}else{
			nextIcon[0] = (<li className="waves-effect" onClick={this.handlePageChange.bind(this, this.state.data.CurrentPage+1)}><a href="#next"><i className="material-icons">chevron_right</i></a></li>)
		}

		return (
			<div className="container" style={{"margin-top": "20px"}}>
				<div className="row valign-wrapper">
					<div className="col s10 valign">
					  <ul className="pagination left">
					    {prevIcon}
						{nodes}
						{nextIcon}
					  </ul>
					</div>
					<div className="col s2 valign">
						<a className="modal-trigger right waves-effect waves-light btn" href="#addFormModal">
							<i className="material-icons left">add</i>Add
						</a>
					</div>
				</div>
			    <div className="row">
					<div className="col s12">
						<FrameworkTable onFrameworkDelete={this.handleFrameworkDelete} data={this.state.data}/>
					</div>
				</div>
				<AddForm onFrameworkSubmit={this.handleAddFrameworkSubmit}/>
			</div>
		);
	}
});

React.render(<Frameworks url="/api/frameworks" />, document.getElementById('container'));

$(document).ready(function(){
	$('.modal-trigger').leanModal();
});
