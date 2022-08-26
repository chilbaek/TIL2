import './App.css';
import * as React from 'react'
import Map from "react-map-gl";

const MAPBOX_TOKEN = 'pk.eyJ1IjoibGpqN2tyIiwiYSI6ImNra3JvcGoxcTB5cG0ybnMxbnk0bDF0ZmUifQ.lEiOIr5N1M6mjDH10aC8MQ';

const changeMapLanguage = (map) => {
    map.getStyle().layers.forEach((layer) => {
        if (layer.id.endsWith('-label')) {
            map.setLayoutProperty(layer.id, 'text-field', [
                'coalesce',
                ['get', 'name_ko'],
                ['get', 'name'],
            ]);
        }
    });

};

function App() {
    return (
        <Map
            initialViewState={{
                longitude: 126.98877,
                latitude: 37.56583,
                zoom: 18
            }}
            style={{width: 600, height: 400}}
            // mapStyle="mapbox://styles/mapbox/streets-v9"
            mapStyle="mapbox://styles/mapbox/light-v10"
            mapboxAccessToken={MAPBOX_TOKEN}
            onStyleLoad={changeMapLanguage}
        />
    );

}

export default App;
