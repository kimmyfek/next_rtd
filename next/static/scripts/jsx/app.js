var Button = require('react-bootstrap/lib/Button');
var ButtonGroup = require('react-bootstrap/lib/ButtonGroup');
var Col = require('react-bootstrap/lib/Col');
var Grid = require('react-bootstrap/lib/Grid');
var PageHeader = require('react-bootstrap/lib/PageHeader');
var Row = require('react-bootstrap/lib/Row');
var Station = require('./Station');
var axios = require ('axios');

var App = React.createClass({
	getInitialState: function(){
		return {
          to: "",
          from: "",
          stations: [],
          reset: false
        };
	},

    startOver: function(me){
      me.setState(this.getInitialState());
      //this.setState({reset:true}, () => {
      //  this.forceUpdate();
      //  });
      this.setState({reset:true});
      axios.get('/stations')
      .then(res => {
        var stationSorted = res.data.sort(function(a, b) {
          var textA = a.name.toUpperCase();
          var textB = b.name.toUpperCase();
          return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
        });
		this.setState({stations:stationSorted});
      });
    },

    componentDidMount: function(){
      axios.get('/stations')
      .then(res => {
        var stationSorted = res.data.sort(function(a, b) {
          var textA = a.name.toUpperCase();
          var textB = b.name.toUpperCase();
          return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
        });
		this.setState({stations:stationSorted});
      });
      this.setState({reset:false});
    },

	render: function(){
	  return (
		<div>
		  <div className="header">
          <div className="reset-btn">
            <Button onClick= {() => this.startOver(this)}
              bsSize="xsmall"
              bsStyle="warning">Start Over
            </Button>
          </div>
          </div>
          <div id="from_stations">
              <Station stations={this.state.stations} reset={this.state.reset} />
          </div>
		</div>
	  );
	}

});


ReactDOM.render( <App />,
  document.getElementById('main')
);

