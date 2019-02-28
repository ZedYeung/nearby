import React from 'react';
import { Tabs, Spin } from 'antd';
import {API_ROOT, AUTH_PREFIX, GEO_OPTIONS, POS_KEY, TOKEN_KEY, GOOGLE_MAP_URL} from "../.env";
import $ from 'jquery';
import {Gallery} from "./Gallery";
import {CreatePostButton} from './CreatePostButton';
import {WrappedAroundMap} from './AroundMap';

const TabPane = Tabs.TabPane;

export class Home extends React.Component{
    state = {
        loadingGeoLocation : false,
        loadingPosts: false,
        error: '',
        posts:[],
    }

    componentDidMount () {
        this.setState({loadingGeoLocation: true, error:'',});
        this.getGeoLocation();
    }

    getGeoLocation(){
        if ("geolocation" in navigator) {
            /* 
            The Navigator.geolocation read-only property returns a Geolocation object that gives Web content access to the location of the device.
            This allows a Web site or app to offer customized results based on the user's location.
            https://developer.mozilla.org/en-US/docs/Web/API/Navigator/geolocation
            */
            navigator.geolocation.getCurrentPosition( // async request
                this.onSuccessLoadGeoLocation,
                this.onFailedLoadGeoLocation,
                GEO_OPTIONS
            );
        } else {
            /* geolocation IS NOT available */
            this.setState({error: 'Your browser does not support geolocation'});
        }
    }

    onSuccessLoadGeoLocation = (position) => {
        console.log(position);
        this.setState({
            loadingGeoLocation: false,
            error:'',
        });
        const {latitude, longitude} = position.coords;
        localStorage.setItem(POS_KEY, JSON.stringify({lat: latitude, lon: longitude}));
        this.loadNearByPosts();
    }

    onFailedLoadGeoLocation =() => {
        this.setState({
            loadingGeoLocation: false,
            error: 'failed to load geo location!'
        });
    }

    getGalleryPanelContent = () => {
        if (this.state.error){
            return <div>{this.state.error}</div>
        } else if (this.state.loadingGeoLocation) {
            // return <span>loading geo location</span>;
            return <Spin tip="Loading geo location..."/>
        } else if (this.state.loadingPosts) {
            return <Spin tip="Loading posts..."/>
        } else if (this.state.posts && this.state.posts.length > 0) {
            const images = this.state.posts.map((post) => {
               return {
                   user: post.user,
                   src: post.url[0],
                   thumbnail: post.url[0],
                   thumbnailWidth: 300,
                   thumbnailHeight: 300,
                   caption: post.message,
               };
            });
            return <Gallery images={images}/>
        }
    }

    loadNearByPosts = (location, radius) => {
        const {lat, lon} = location? location: JSON.parse(localStorage.getItem(POS_KEY));
        // const {lat, lon} = {lat:47.7915953, lon:-122.3937977};
        const range = radius ? radius : 20;

        this.setState({loadingPosts:true, error:''});
        return $.ajax({
            url: `${API_ROOT}/search?lat=${lat}&lon=${lon}&range=${range}`,
            method: 'GET',
            headers: {
                Authorization: `${AUTH_PREFIX} ${localStorage.getItem(TOKEN_KEY)}`
            },
        }).then((response) => {
            this.setState({posts:response, loadingPosts:false, error:''});
            console.log("posts: " + response);
        }, (error) => {
            this.setState({loadingPosts:false, error:error.responseText});
        }).catch((error) => {
            console.log(error);
        });
    }

    render(){
        const createPostButton = <CreatePostButton loadNearByPosts={this.loadNearByPosts}/>;

        return (//jsx ==  React.createElement(..)
            <Tabs tabBarExtraContent={createPostButton} className="main-tabs">
                <TabPane tab="Posts" key="1">
                    {this.getGalleryPanelContent()}
                </TabPane>
                <TabPane tab="Map" key="2">
                    <WrappedAroundMap
                        googleMapURL={GOOGLE_MAP_URL}
                        loadingElement={<div style={{ height: `100%` }} />}
                        containerElement={<div style={{ height: `600px` }} />}
                        mapElement={<div style={{ height: `100%` }}/>}
                        posts = {this.state.posts}
                        loadNearByPosts = {this.loadNearByPosts}
                    />
                </TabPane>
            </Tabs>
        );
    }
}