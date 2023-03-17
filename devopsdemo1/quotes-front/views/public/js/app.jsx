class App extends React.Component {
  render() {
 //   if (this.loggedIn) {
      return (<LoggedIn />);
 //   } else {
 //     return (<LoggedOut />);
 //   }
  }
}

class LoggedOut extends React.Component {
  render() {
    return (
      <div className="container">
        <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
          <h1>Hello World - Logged out</h1>
          <p>Une nouvelle App</p>
          <p>Sign in to get access </p>
          <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
        </div>
      </div>
    )
  }
}

class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      citations: []
    }

    this.serverRequest = this.serverRequest.bind(this);
    this.logout = this.logout.bind(this);

   }

    logout() {
      localStorage.removeItem("id_token");
      localStorage.removeItem("access_token");
      localStorage.removeItem("profile");
      location.reload();
    }
    
    refresh() {
      location.reload();
    }

    serverRequest() {
 //     $.get("http://34.98.114.247/api/citations", res => {
      $.get("https://app.dev.gabrielbechara.com/api/citations", res => {
        this.setState({
          citations: res
        });
      });

  }
  
  componentDidMount() {
    this.serverRequest();
  }

  render() {
    return (
      <div className="container">
        <div className="col-xs-8 col-xs-offset-2 jumbotron">
            <div className="col-lg-12">
              <br />
              <span className="pull-right"><a onClick={this.refresh}>Refresh</a> <a onClick={this.logout}>Log out</a></span>
              <h2>Hello World v.0.3</h2>
              <div className="row">
                {this.state.citations.map(function(citation, i){
                  return (<Citation key={i} citation={citation} />);
                })}
              </div>
            </div>
        </div>
      </div>
    )
  }
}

class Citation extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: "",
      citations: [],
    };
    this.like = this.like.bind(this);
    this.serverRequest = this.serverRequest.bind(this);
  }
    
  like() {
    let citation = this.props.citation;
    this.serverRequest(citation);
  }

  refresh() {
    location.reload();
  }
  
  serverRequest(citation) {
    $.post(
 //     "http://34.98.114.247/api/citations/like/" + citation.id,
      "https://app.dev.gabrielbechara.com/api/citations/like/" + citation.id,
      { like: 1 },
      res => {
        console.log("res... ", res);
        this.setState({ liked: "Liked!", citation: res });
        this.props.citation = res;
      }
    );
    this.refresh()
  }

  
  render() {
    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading" style={{backgroundColor: this.props.citation.color, color: '#ffffff', fontWeight: 'bold'}}>#{this.props.citation.id} <span className="pull-right">{this.state.liked}</span></div>
          <div className="panel-body" style={{backgroundColor: this.props.citation.color, color: '#ffffff', fontWeight: 'bold'}}>
            {this.props.citation.citation}
          </div>
          
          <div className="panel-footer" style={{backgroundColor: this.props.citation.color, color: '#ffffff', fontWeight: 'bold'}}>
            {this.props.citation.likes} Likes &nbsp;
            <a onClick={this.like} className="btn btn-default">
              <span className="glyphicon glyphicon-thumbs-up"></span>
            </a>
          </div>
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('monapp'));
