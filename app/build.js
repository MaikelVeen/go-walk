import esbuild from 'esbuild';
import path from 'path';
import fs from 'fs';
import dotenv from 'dotenv';

dotenv.config();

const entryPoints = ['src/app.ts'];
const outdir = 'dist';
const dataDir = 'data';
const manifestPath = path.join(dataDir, 'manifest.json');

const isWatchMode = process.argv.includes('--watch');

let geojsonDataMap = {};
try {
  if (fs.existsSync(manifestPath)) {
    const manifestContent = fs.readFileSync(manifestPath, 'utf-8');
    const filenames = JSON.parse(manifestContent);
    
    if (Array.isArray(filenames)) {
      console.log(`Reading ${filenames.length} GeoJSON files listed in manifest: ${manifestPath}`);
      
      filenames.forEach(filename => {
        const filePath = path.join(dataDir, filename);
        
        try {
          if (fs.existsSync(filePath)) {
            const fileContent = fs.readFileSync(filePath, 'utf-8');
            geojsonDataMap[filename] = JSON.parse(fileContent);
            console.log(`  - Embedded data for ${filename}`);
          } else {
            console.warn(`  - Warning: GeoJSON file not found: ${filePath}`);
          }
        } catch (fileError) {
          console.error(`  - Error reading or parsing GeoJSON file ${filePath}:`, fileError);
        }
      });
    } else {
      console.warn(`Warning: Invalid manifest format in ${manifestPath}. Expected an array.`);
    }
  } else {
    console.warn(`Warning: Manifest file not found at ${manifestPath}. No GeoJSON data will be injected.`);
  }
} catch (error) {
  console.error(`Error reading or parsing manifest file ${manifestPath}:`, error);
}

const mapboxToken = process.env.MAPBOX_ACCESS_TOKEN || '';
if (!mapboxToken) {
    console.warn("Warning: MAPBOX_ACCESS_TOKEN environment variable not set or empty. Map might not load correctly.");
}

esbuild.build({
  entryPoints: entryPoints,
  bundle: true,
  outfile: path.join(outdir, 'app.js'),
  platform: 'browser',
  sourcemap: true,
  define: {
    '__GEOJSON_DATA__': JSON.stringify(geojsonDataMap),
    'process.env.MAPBOX_ACCESS_TOKEN': JSON.stringify(mapboxToken)
  },
  ...(isWatchMode && {
      watch: {
          onRebuild(error, result) {
              if (error) console.error('watch build failed:', error);
              else console.log('watch build succeeded:', result);
          }
      }
  }),
  logLevel: 'info',
}).then(result => {
    if (isWatchMode) {
        console.log('watching...');
    } else {
        console.log('Build finished.');
    }
}).catch(() => process.exit(1)); 