import React from 'react';
import { Marker, InfoWindow} from 'react-google-maps';
export class AroundMarker extends React.Component {
    state = {
        isOpen : false,
    }
    onToggleOpen = () =>{
        this.setState((prevState)=>({isOpen:!prevState.isOpen}));
    }

    render(){
        const {user, message, location, url} = this.props.post;
        return (
            <Marker
                position={{ lat: location.lat, lng: location.lon }}
                onMouseOver={this.onToggleOpen}
                onMouseOut={this.onToggleOpen}
            >
                {this.state.isOpen ?
                    <InfoWindow onCloseClick={this.onToggleOpen}>
                        <div>
                            <img className="around-marker-image" src={url[0]} alt={`${user}: ${message}`}/>
                            <p>{`${user}: ${message}`}</p>
                        </div>
                    </InfoWindow> : null}
            </Marker>
        );
    }
}