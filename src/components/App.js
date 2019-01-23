import React from 'react';
import '../styles/App.css';
import {Header} from './Header';
import  {Main} from './main';
import {TOKEN_KEY} from '../.env'

class App extends React.Component {
  state = {
      isLoggedIn: !! localStorage.getItem(TOKEN_KEY),
  }

  handleLogin = (response) => {
      localStorage.setItem(TOKEN_KEY, response);
      this.setState({isLoggedIn:true});

  }

  handleLogout = () =>{
      localStorage.removeItem(TOKEN_KEY);
      this.setState({isLoggedIn: false});
  }

  render() {
    return (
      <div className="App">
        <Header isLoggedIn={this.state.isLoggedIn} handleLogout={this.handleLogout}/>
        <Main isLoggedIn={this.state.isLoggedIn} handleLogin={this.handleLogin} />
      </div>
    );
  }
}

export default App;