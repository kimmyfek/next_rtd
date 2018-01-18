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
          stations: [],
          reset: false
        };
        this.state = this.initState;
    }

    startOver(me){
      me.setState(this.initState);
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
      this.setState({reset:false});
    }

	render(){
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
          <br />
		</div>
	  );
	}

}

ReactDOM.render( <App />,
  document.getElementById('main')
);
