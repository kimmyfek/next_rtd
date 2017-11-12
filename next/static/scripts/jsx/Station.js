var React = require('react');
var Button = require('react-bootstrap/lib/Button');
var ButtonGroup = require('react-bootstrap/lib/ButtonGroup');
var Col = require('react-bootstrap/lib/Col');
var Grid = require('react-bootstrap/lib/Grid');
var PageHeader = require('react-bootstrap/lib/PageHeader');
var Row = require('react-bootstrap/lib/Row');
var Time = require('./Time');

var Station = React.createClass({
 	// sets initial state
	getInitialState: function(){
		return {
          stations: this.props.stations,
          isDestination: false,
          to: "",
          from: "",
          times: "",
          direction: "",
          pageText: "Departing Station"
        };
	},

	getConnectingStations: function(station){
        if(this.state.isDestination == false){
          this.setState({stations: station.connections});
          this.setState({isDestination: true});
          this.setState({from: station.name});
          this.setState({pageText: "Arriving Station"});
        } else {
          this.setState({to: station.name});
        }
	},


    getNextTimes: function(to, from){
      this.setState({times: [{"time": "2:00pm", "line": "C"}, {"time": "2:01pm", "line": "C"}, {"time": "2:02pm", "line": "E"}]});
      this.setState({direction: "Northbound"});
    },

    shouldComponentUpdate: function(nextProps, nextState){
      if(this.state.times == ""){
        return true;
      } else{
        return false;
      }
    },

 	render: function() {
		var me = this;

        if(me.state.to != "" && me.state.from != ""){
            // we have our to & from  stations,
            // lets call the backend for the times
            me.getNextTimes(me.state.to, me.state.from);

            if(me.state.times != ""){
              return (<Time to={me.state.to}
                        from={me.state.from}
                        times={me.state.times}
                        direction={me.state.direction} /> );
            } else {
              return (<div>Loading</div>);
            }
        } else {

		var listStations = this.state.stations.map(function(station) {
			return (
			  <Button
                key={station.name}
				onClick={ () => {me.getConnectingStations(station)}}>
			   {station.name}
			  </Button>
			  );
		  });
      return (
        <div>
        <h3>{me.state.pageText}</h3>
        <br />
        <ButtonGroup vertical block>
          {listStations}
        </ButtonGroup>
        </div>
      );

      }


 	}
});

module.exports = Station;
