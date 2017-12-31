var React = require('react');
var Button = require('react-bootstrap/lib/Button');
var ButtonGroup = require('react-bootstrap/lib/ButtonGroup');
var Col = require('react-bootstrap/lib/Col');
var Grid = require('react-bootstrap/lib/Grid');
var PageHeader = require('react-bootstrap/lib/PageHeader');
var Row = require('react-bootstrap/lib/Row');
var Time = require('./Time');
var axios = require('axios');

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
          // sort the station names
          var conn = station.connections.sort(function(a, b) {
			  var textA = a.name.toUpperCase();
			  var textB = b.name.toUpperCase();
			  return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
		  });

          this.setState({stations: station.connections});
          this.setState({isDestination: true});
          this.setState({from: station.name});
          this.setState({pageText: "Arriving Station"});
        } else {
          this.setState({to: station.name});
        }
	},


    getNextTimes: function(to, from){
	  to = encodeURIComponent(to);
	  from = encodeURIComponent(from);
	  axios.get('/times?to=' + to + '&from=' + from)
      .then(res => {
        // times should already be sorted; commenting out
        //  var times= res.data.sort(function(a, b) {
		//	  var textA = a.arrival_time;
		//	  var textB = b.arrival_time;
		//	  return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
		//  });
        this.setState({times:res.data});
      });
    },

    componentWillReceiveProps(nextProps){
      if(nextProps.reset){
          this.setState(this.getInitialState());
      }
      this.setState({stations:nextProps.stations});
      this.forceUpdate();
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
              return (<div className="loading"></div>);
            }
        } else {

		var listStations = this.state.stations.map(function(station) {
			return (
            <div className="station" key={"station" + station.name}>
			  <Button
                key={"btn"+station.name}
				onClick={ () => {me.getConnectingStations(station)}}>
                <div key={station.name} className="station-btn-text">
			      {station.name}
               </div>
			  </Button>
              <br />
            </div>
			  );
		  });
      return (
        <div>
          <h3>{me.state.pageText}</h3>
        <br />
        <ButtonGroup vertical block className="stations">
          {listStations}
        </ButtonGroup>
        </div>
      );

      }


 	}
});

module.exports = Station;
