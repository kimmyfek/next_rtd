var Button = require('react-bootstrap/lib/Button');
var ButtonGroup = require('react-bootstrap/lib/ButtonGroup');
var Col = require('react-bootstrap/lib/Col');
var Grid = require('react-bootstrap/lib/Grid');
var PageHeader = require('react-bootstrap/lib/PageHeader');
var Row = require('react-bootstrap/lib/Row');
//const stationStyles = {maxWidth: 400, margin: '0 auto 10px'};

var StationSelect = React.createClass({

  getInitialState: function() {
    return {
      isActive: false
    };
  },

  render: function() {
    let isActive = this.state.isActive;
    return (
      <div>
        <PageHeader>Next RTD</PageHeader>
        <Grid>
          <Row className="show-grid">
            <Col md={4}>
            </Col>
            <Col md={4}>
            <div id="from_stations">
              <ButtonGroup vertical block>
                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Union Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Pepsi Center - Elitch Gardens Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Sports Authority Field at Mile High Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Auraria West Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  10th and Osage Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Alameda Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  I25 & Broadway Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Evans Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Englewood Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Oxford Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Littleton Downtown Station
                </Button>

                <Button
                  block
                  onClick={this.thing}
                  active={isActive}>
                  Littleton Mineral Station
                </Button>

              </ButtonGroup>
            </div>
            </Col>
            <Col md={4}>
            </Col>
          </Row>
        </Grid>
      </div>
    )
  },

  thing: function() {
    this.setState({isActive: true});
  }

});

function Stations(props) {

}

const stations = [
  {"name": "Union Station", "isActive": false},
  {"name": "Pepsi Center - Elitch Gardens Station", "isActive": false},
  {"name": "Sports Authority Field at Mile High Station", "isActive": false},
  {"name": "Auraria West Station", "isActive": false},
  {"name": "10th and Osage Station", "isActive": false},
  {"name": "Alameda Station", "isActive": false},
  {"name": "I25 & Broadway Station", "isActive": false},
  {"name": "Evans Station", "isActive": false},
  {"name": "Englewood Station", "isActive": false},
  {"name": "Oxford Station", "isActive": false},
  {"name": "Littleton Downtown Station", "isActive": false},
  {"name": "Littleton Mineral Station", "isActive": false}
]

ReactDOM.render(
  <Stations stations={stations}>
  document.getElementById('main')
);
