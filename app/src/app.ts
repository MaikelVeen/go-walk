import mapboxgl from 'mapbox-gl';

declare const __GEOJSON_DATA__: { [filename: string]: any };

interface MapOptions {
  container: string;
  style: string;
  center: [number, number];
  zoom: number;
}

class MapboxMap {
  private map: mapboxgl.Map;

  constructor(apiToken: string, options: MapOptions) {
    mapboxgl.accessToken = apiToken;

    this.map = new mapboxgl.Map({
      container: options.container,
      style: options.style,
      center: options.center,
      zoom: options.zoom
    });

    this.map.addControl(new mapboxgl.NavigationControl(), 'top-right');

    this.map.addControl(
      new mapboxgl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true
        },
        trackUserLocation: true
      })
    );

    this.setupEventListeners();
  }

  private setupEventListeners(): void {
    this.map.on('load', () => {
      try {
        const geojsonDataMap = __GEOJSON_DATA__;

        if (typeof geojsonDataMap !== 'object' || geojsonDataMap === null) {
          console.error("__GEOJSON_DATA__ is not an object. Build process might have failed to inject data.");
          return;
        }

        const filenames = Object.keys(geojsonDataMap);

        if (filenames.length === 0) {
            console.log("No GeoJSON data found in injected bundle. Nothing to load.");
            return;
        }

        console.log(`Found ${filenames.length} datasets in injected data. Adding to map...`);

        filenames.forEach(filename => {
          const geojsonData = geojsonDataMap[filename];
          const id = `walk-path-${filename.split('.')[0] ?? Math.random().toString(36).substring(7)}`;

          if (this.map.getSource(id)) {
            console.warn(`Source with ID ${id} already exists. Skipping.`);
            return;
          }

          try {
            this.map.addSource(id, {
              type: 'geojson',
              data: geojsonData
            });

            this.map.addLayer({
              id: `${id}-layer`,
              type: 'line',
              source: id,
              layout: {
                'line-join': 'round',
                'line-cap': 'round'
              },
              paint: {
                'line-color': '#008100',
                'line-width': 3,
                'line-opacity': 0.8
              }
            });
            console.log(`Successfully added embedded GeoJSON layer: ${id}`);
          } catch(mapError) {
            console.error(`Error adding source/layer for ${filename} (ID: ${id}):`, mapError);
          }
        });
      } catch (error) {
        console.error("Failed to process embedded GeoJSON data:", error);
      }
    });

    this.map.on('click', (e) => {
      console.log(`Map clicked at: [${e.lngLat.lng}, ${e.lngLat.lat}]`);
    });
  }

  public getCenter(): [number, number] {
    const center = this.map.getCenter();
    return [center.lng, center.lat];
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const accessToken = process.env.MAPBOX_ACCESS_TOKEN;

  if (!accessToken) {
    console.error("Mapbox Access Token not configured! Set MAPBOX_ACCESS_TOKEN in your .env file and rebuild.");
    const mapDiv = document.getElementById('map');
    if (mapDiv) {
        mapDiv.innerHTML = '<div style="padding: 20px; text-align: center; font-family: sans-serif; color: red;">Mapbox Access Token is missing. Please check configuration.</div>';
    }
    return;
  }

  const mapOptions: MapOptions = {
    container: 'map',
    style: 'mapbox://styles/mapbox/light-v10',
    center: [4.4777, 51.9244],
    zoom: 12
  };

  const mapManager = new MapboxMap(accessToken, mapOptions);
});

export default MapboxMap;