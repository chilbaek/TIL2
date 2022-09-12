import './App.css';
import * as React from 'react'
import Map, {Layer, Source} from "react-map-gl";

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
  const mapGL = this.getMap();
  const bounds = mapGL.getBounds();
  console.log(bounds);
  const geojson = {
    type: 'Feature',
    geometry: {
      type: 'MultiLineString',
      coordinates: [
        [
          [
            127.06847450173576,
            37.53525157645467
          ],
          [
            127.0717755195668,
            37.53429708023673
          ],
          [
            127.07567058614043,
            37.53315233197593
          ]
        ]
      ]
    },
    properties: {
      'CONNECT': '000',
      'F_NODE': '1040002900',
      'LANES': 2,
      'LENGTH': 677.361,
      'LINK_ID': '1040006600',
      'MAX_SPD': 30,
      'MULTI_LINK': '0',
      'ROAD_NAME': '뚝섬로',
      'ROAD_NO': '-',
      'ROAD_RANK': '104',
      'ROAD_TYPE': '000',
      'ROAD_USE': '0',
      'T_NODE': '1040002700'
    }
  };

  const layerStyle = {
    id: 'route',
    type: 'line',
    paint: {
      'line-width': 5,
      'line-color': '#007cbf'
    }
  };

  return (
    <Map
      initialViewState={{
        longitude: 126.98877, latitude: 37.56583,
        zoom: 18
      }} 
      style={{width: 600, height: 400}}
      // mapStyle="mapbox://styles/mapbox/streets-v9" 
      mapStyle="mapbox://styles/mapbox/light-v10" 
      mapboxAccessToken={MAPBOX_TOKEN} Ï
      onStyleLoad={changeMapLanguage}
    >
      <Source id="my-data" type="geojson" data={geojson}>
        <Layer {...layerStyle} />
      </Source>
    </Map>
  );
}

export default App;
