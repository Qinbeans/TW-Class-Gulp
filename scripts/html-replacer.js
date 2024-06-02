// used in require('./scripts/html-replacer')({ classMap })

const { Transform } = require('stream');
const { JSDOM } = require('jsdom');
const { parse } = require('node-html-parser');

const classMap = require('../classMap.json');
const { GCProfiler } = require('v8');

module.exports = function (options) {
    const stream = new Transform({ objectMode: true });
    stream._transform = function (file, encoding, callback) {
        const dom = new JSDOM(file.contents.toString());
        const document = dom.window.document;
        const root = parse(file.contents.toString());
    
        root.querySelectorAll('*').forEach(node => {
            if (node.classNames.length > 0) {
                const classNames = node.classNames.split(' ').sort().join(' ');
                if (classMap[classNames]){
                    // clear all classes
                    classNames.split(' ').forEach(className => {
                        node.classList.remove(className);
                    });
                    node.classList.add(classMap[classNames]);
                }
            }
        });
    
        file.contents = Buffer.from(root.toString());
        this.push(file);
        callback();
    };
    
    return stream;
}
