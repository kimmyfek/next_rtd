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
			return (
			  <h4 key={time.time + time.line}>
              <Well>
			   Line {time.line} @ {time.time}
			  </Well></h4>
			  );
		  });

      return (
        <div>
        <h3><span>{me.state.from} to {me.state.to} </span></h3>
        <Label bsStyle="primary">{me.state.direction} </Label>
        <br />
        <br />
        <br />
        <ListGroup vertical block>
          {listTimes}
        </ListGroup>
        </div>
      );

 	}
});

module.exports = Time;
