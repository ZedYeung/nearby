import React, { Component } from 'react';
import { Register } from "./Register";
import { Login } from "./Login";
import { Switch, Route, Redirect } from 'react-router-dom';
import { Home } from './Home';

export class Main extends Component{
    getLogin = () => {
        return this.props.isLoggedIn? <Redirect to="/home"/> : <Login handleLogin = {this.props.handleLogin}/>;
    }

    getHome = () => {
        return this.props.isLoggedIn? <Home setLocation={this.props.setLocation} getLocation={this.props.getLocation}/> : <Redirect to="/login"/>;
    }

    getRoot = () => {
        return <Redirect to ='/login'/>
    }

    render(){
        return (
            <section className="main">
                <Switch>
                    <Route exact path="/" render={this.getRoot}/>
                    <Route path="/register" component={Register}/>
                    <Route path="/login" render={this.getLogin}/>
                    <Route path="/home" render={this.getHome}/>
                    <Route render={this.getRoot}/>
                </Switch>
            </section>
        )
    }
}