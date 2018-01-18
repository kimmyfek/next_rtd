import React from 'react';
import { Button, ButtonGroup, Col, Grid, PageHeader, Row } from 'react-bootstrap';
import axios from 'axios';
import Time from './Time';

class Station extends React.Component {
 	// sets initial state
    constructor(props) {
        super(props);
        this.initState = {stations: this.props.stations,
          isDestination: false,
          to: "",
          from: "",
          times: "",
          direction: "",
          pageText: "Where Are You Leaving From?"
        };
        this.state = this.initState;
	}

	getConnectingStations(station){
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
          this.setState({pageText: "Where Are You Going?"});
        } else {
          this.setState({to: station.name});
        }
	}


    getNextTimes(to, from){
	  to = encodeURIComponent(to);
	  from = encodeURIComponent(from);
	  axios.get('/times?to=' + to + '&from=' + from)
      .then(res => {
        this.setState({times:res.data});
      });
    }

    componentWillReceiveProps(nextProps){
      this.setState({stations:nextProps.stations});
      this.forceUpdate();
    }

    shouldComponentUpdate(nextProps, nextState){
      if(this.state.times == ""){
        return true;
      } else{
        return false;
      }
    }

 	render() {
		var me = this;
        if(me.state.to != "" && me.state.from != ""){
            // we have our to & from  stations,
            // lets call the backend for the times
            if(me.state.times == ""){
              me.getNextTimes(me.state.to, me.state.from);
            }

            if(me.state.times != ""){
              return (
                <Time to={me.state.to}
                from={me.state.from}
                times={me.state.times}
                direction={me.state.direction} />
                );
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

          if(me.state.from != ""){
              var fromDisplay = "from: " + me.state.from;
          } else {
              var fromDisplay = "";
          }
          return (
            <div>
            <div className="header">
              <h3>{me.state.pageText}</h3>
              <h5><i>{fromDisplay}</i></h5>
            </div>
            <br />
            <ButtonGroup vertical block className="stations">
              {listStations}
            </ButtonGroup>
            </div>
          );
        }
 	}
}

export default Station;
