import React from 'react';
import ReactDOM from 'react-dom';
import { Button, ButtonGroup, Col, Grid, PageHeader, Row } from 'react-bootstrap';
import axios from 'axios';
import Station from './Station';

class App extends React.Component {
	constructor(props) {
        super(props);
        this.initState = {to: "",
          from: "",
          stations: []
        };
        this.state = this.initState;
    }

    startOver(me){
      me.setState(this.initState);
      axios.get('/stations')
      .then(res => {
        var stationSorted = res.data.sort(function(a, b) {
          var textA = a.name.toUpperCase();
          var textB = b.name.toUpperCase();
          return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
        });
		this.setState({stations:stationSorted});
      });
    }

    componentDidMount(){
      axios.get('/stations')
      .then(res => {
        var stationSorted = res.data.sort(function(a, b) {
          var textA = a.name.toUpperCase();
          var textB = b.name.toUpperCase();
          return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
        });
		this.setState({stations:stationSorted});
      });
    }

	render(){
	  return (
		<div>
          <div id="from_stations">
              <Station stations={this.state.stations} />
          </div>
          <br />
		</div>
	  );
	}

}

ReactDOM.render( <App />,
  document.getElementById('main')
);
