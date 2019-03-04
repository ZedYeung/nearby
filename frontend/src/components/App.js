import React from 'react';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import '../styles/App.css';
import { Header } from './Header';
import { Main } from './main';

class App extends React.Component {
  static propTypes = {
    cookies: instanceOf(Cookies).isRequired
  };

  state = {
    isLoggedIn: !! this.props.cookies.get(process.env.REACT_APP_TOKEN_KEY),
  }

  handleLogin = (res) => {
      console.log(res.data)
      this.props.cookies.set(process.env.REACT_APP_TOKEN_KEY, res.data, { path: '/', maxAge: 3600 * 24 });
      this.setState({isLoggedIn:true});
  }

  handleLogout = () =>{
    this.props.cookies.remove(process.env.REACT_APP_TOKEN_KEY);
    this.setState({isLoggedIn: false});
  }

  setLocation = (latitude, longitude) => {
    this.props.cookies.set(process.env.REACT_APP_POS_KEY, JSON.stringify({
      lat: latitude,
      lon: longitude
    }), {
      path: '/',
      maxAge: 3600 * 24
    });
  }

  getLocation = () => {
    return this.props.cookies.get(process.env.REACT_APP_POS_KEY)
  }

  render() {
    return (
      <div className="App">
        <Header isLoggedIn={this.state.isLoggedIn} handleLogout={this.handleLogout}/>
        <Main
          isLoggedIn={this.state.isLoggedIn}
          handleLogin={this.handleLogin}
          setLocation={this.setLocation}
          getLocation={this.getLocation}
        />
      </div>
    );
  }
}

export default withCookies(App);