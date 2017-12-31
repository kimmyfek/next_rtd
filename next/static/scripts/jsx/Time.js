var React = require('react');
var Button = require('react-bootstrap/lib/Button');
var Well = require('react-bootstrap/lib/Well');
var ButtonGroup = require('react-bootstrap/lib/ButtonGroup');
var ListGroup = require('react-bootstrap/lib/ListGroup');
var ListGroupItem = require('react-bootstrap/lib/ListGroupItem');
var Label = require('react-bootstrap/lib/Label');
var Col = require('react-bootstrap/lib/Col');
var Grid = require('react-bootstrap/lib/Grid');
var PageHeader = require('react-bootstrap/lib/PageHeader');
var Row = require('react-bootstrap/lib/Row');

var Time = React.createClass({
 	// sets initial state
	getInitialState: function(){
		return {
          to: this.props.to,
          from: this.props.from,
          times: this.props.times,
          direction: this.props.direction
        };
	},


 	render: function() {
		var me = this;
		var listTimes = me.state.times.map(function(time) {
            // format the time
            var d_time = time.departure_time;
            var d_hr = parseInt(d_time.substring(0,2));
            var d_min = d_time.substring(3,5);
            var am_pm = 'AM';
            if (d_hr >= 24){
                d_hr -= 24;
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
        <ListGroup vertical block className="time-list">
          {listTimes}
        </ListGroup>
        </div>
      );

 	}
});

module.exports = Time;
