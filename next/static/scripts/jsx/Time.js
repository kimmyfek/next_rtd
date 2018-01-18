import React from 'react';
import { Button, ButtonGroup, ListGroup, ListGroupItem, Label, Well, Col, Grid, PageHeader, Row } from 'react-bootstrap';

class Time extends React.Component {
 	// sets initial state
    constructor(props) {
        super(props);
        this.state = {to: this.props.to,
          from: this.props.from,
          times: this.props.times,
          direction: this.props.direction
        };
	}

    formatTime(t) {
      // format the time
      var d_time = t;
      var d_hr = parseInt(d_time.substring(0,2));
      var d_min = d_time.substring(3,5);
      var am_pm = 'AM';
      if (d_hr >= 24){
          d_hr -= 24;
          if (d_hr == 0){
              d_hr = 12;
          }
          am_pm = 'AM';
      } else if (d_hr >=12 && d_hr < 24){
          if (d_hr != 12){
              d_hr -= 12;
          }
          am_pm = 'PM';
      }
      return d_hr.toString() + ":" + d_min.toString() + " " + am_pm;
    }

 	render() {
		var me = this;

        if (me.state.times == null){
            return (
              <div>
              <div className="header">
              <h3><span>{me.state.from} to {me.state.to} </span></h3>
              </div>
              <hr />
              <br />
              <div className="alert alert-warning">
              <span>There are no available times right now </span>
              </div>
              </div>
            );
        } else {
            var listTimes = me.state.times.map(function(time) {
                var dep_time = me.formatTime(time.departure_time);
                var arr_time = me.formatTime(time.arrival_time);

                return (
                  <div key={time.route+dep_time}>
                  <ListGroupItem header={time.route + " line at " + dep_time} bsStyle="success">
                    <i>{"You should arrive by " + arr_time}</i>
                  </ListGroupItem>
                  <br />
                  </div>
                  );
              });

              return (
                <div>
                <div className="header">
                <h3><span>{me.state.from} to {me.state.to} </span></h3>
                </div>
                <hr />
                <br />
                <ListGroup className="time-list">
                  {listTimes}
                </ListGroup>
                </div>
              );
        }
 	}
}

export default Time;
