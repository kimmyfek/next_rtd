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
                // format the time
                var d_time = time.departure_time;
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

                return (
                  <h4 key={time.departure_time + time.route}>
                  <Well>
                    <span className="next-train">
                    {time.route}
                    </span>
                    <span className="next-time">
                    {d_hr.toString() + ":" + d_min.toString() + " " + am_pm}
                    </span>
                  </Well></h4>
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
