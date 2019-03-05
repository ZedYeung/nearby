import React, { Component } from 'react';
import { Tabs, Spin } from 'antd';
import { Gallery } from "./Gallery";
import { CreatePostButton } from './CreatePostButton';
import { WrappedAroundMap } from './AroundMap';
import { getSearch } from './API';

const TabPane = Tabs.TabPane;
const GEO_OPTIONS = {
    enableHighAccuracy: false,
    maximumAge        : 3600000, // milliseconds a possible cached position to return. set 0 means not use cache, always try to retrieve real current position
    timeout           : 27000  // milliseconds the device is allowed to take in order to return a position
};

export class Home extends Component{
    state = {
        loadingGeoLocation : false,
        loadingPosts: false,
        error: '',
        posts:null,
    }

    componentDidMount () {
        this.setState({
          loadingGeoLocation: true,
          error:''
        }, () => {
          this.getGeoLocation();
        });
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
        // console.log(position);
        this.setState({
            loadingGeoLocation: false,
            error:'',
        });
        const {latitude, longitude} = position.coords;
        this.props.setLocation(latitude, longitude);
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
        const {lat, lon} = location? location: this.props.getLocation();
        const range = radius ? radius : 20;

        this.setState({
          loadingPosts:true,
          error:''
        });

        return getSearch({
            lat: lat,
            lon: lon,
            range: range
        }).then((res) => {
            this.setState({
              posts:res.data,
              loadingPosts:false,
              error:''
            });
            // console.log("posts: " + res.data);
        }).catch((err) => {
            console.log(err);
        })
    }

    render(){
        const createPostButton = <CreatePostButton getLocation={this.props.getLocation} loadNearByPosts={this.loadNearByPosts}/>;

        return (//jsx ==  React.createElement(..)
            <Tabs tabBarExtraContent={createPostButton} className="main-tabs">
                <TabPane tab="Posts" key="1">
                    {this.getGalleryPanelContent()}
                </TabPane>
                <TabPane tab="Map" key="2">
                    <WrappedAroundMap
                        googleMapURL={process.env.REACT_APP_GOOGLE_MAP_URL}
                        loadingElement={<div style={{ height: `100%` }} />}
                        containerElement={<div style={{ height: `600px` }} />}
                        mapElement={<div style={{ height: `100%` }}/>}
                        posts = {this.state.posts}
                        loadNearByPosts = {this.loadNearByPosts}
                        getLocation={this.props.getLocation}
                    />
                </TabPane>
            </Tabs>
        );
    }
}