import React, { Component } from 'react';
import logo from '../assets/images/logo.svg';
import {Icon} from 'antd';
import PropTypes from 'prop-types';

export class Header extends Component{
    static propTypes = {
        isLoggedIn : PropTypes.bool.isRequired,
        handleLogout: PropTypes.func.isRequired,
    }

    render(){
        return (
            <header className="App-header">
                <img src={logo} className="App-logo" alt="logo" />
                <h1 className="App-title">Nearby</h1>
                {
                    this.props.isLoggedIn &&
                        <a href=""
                           className="logout"
                            onClick={this.props.handleLogout}>
                            <Icon type="logout"/>{' '}logout
                        </a>
                }
            </header>
        )
    }
}