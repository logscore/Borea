// sessionTrack.js

const Borea = {};

{ // open file block

// This scopes all variables to this block

const DOMAIN = "";
const metadataKey = 'metadata';
const postData = true;
const hostname = window.location.hostname;
const GO_PORT = '8080';
const postSessionDataRoute = 'postSession';

Borea.init = function () {
    this[metadataKey] = {
        token: null,
        userId: null,
        sessionId: this.helpers.generateUUID(),
        // previousSessionId: this.getPreviousSessionId(),
        lastActivityTime: this.getLastActivityTime(),
        startTime: new Date(),
        sessionDuration: null,
        userAgent: navigator.userAgent,
        // screenResolution: [`${window.screen.width}x${window.screen.height}`],
        // location set via helper function below
        location: null,
        language: navigator.language,
        referrer: this.helpers.getReferrer(document.referrer),
        // userPath: [],
        // to store and get access to at anytime. these are props you want to be tracked with an event
        // customProperties: {},
    };
    this.options = {
        events: {
            global: {
                options: {
                    /*
                    target,
                    listener,
                    options: {},
                    useCapture,
                    */
                },
            },
            single: {
                type: {
                    options: {
                        /*
                        target,
                        selector,
                        listener,
                        options: {},
                        useCapture,
                        */
                    },
                },
            }
        },
    };
    this.events = {};
    this.eventsArray = []; // TODO determine if adding event obj to array for event order is helpful
    this.enabledEventTypes = null;
    this.defaultEventCallback = null;

    const metadata = sessionStorage.getItem(metadataKey);
    if (metadata != null) {
        // this[metadataKey] = Object.assign(this[metadataKey], JSON.parse(metadata));
        this[metadataKey] = JSON.parse(metadata);
    }

    // tmp, location via ip address
    // this.metadata.location = this.helpers.fetchIPAddress();
    this.initMaintenanceEventListeners();
};

// idk if this is the right idea yet...
Borea.parseOptions = function () {
    if (this.options.userId != undefined) {
        this.setUserId(this.options.userId);
    }

    // if (this.options.customProperties != undefined) {
    //     for (const [key, value] of Object.entries(this.options.customProperties)) {
    //         this.setCustomProperty(key, value);
    //     }
    // }

    // if (this.options.customEventListenerOptions != undefined) {
        // determine what this obj will include...
    // }
};

// Metadata and Options Getters and Setters (block for collapsing)
{
    Borea.updateOptions = function (options) {
        if (typeof options === 'object' && options !== null) {
            for (const [key, value] of Object.entries(options)) {
                this.options[key] = value;
            }
            this.parseOptions();
        } else {
            console.error('Invalid options');
        }
    };

    Borea.setOptions = function (options) {
        if (typeof options === 'object' && options !== null) {
            this.options = options;
            this.parseOptions();
        } else {
            console.error('Invalid options');
        }
    };

    Borea.getUserId = function () {
        return this[metadataKey].userId;
    };

    Borea.setUserId = function (value) {
        if (typeof value === 'string' && value.trim() !== '') {
            this[metadataKey].userId = value.trim();
            localStorage.setItem('userId', value.trim()); // Save to localStorage
        } else {
            console.error('Invalid user ID');
        }
    };

    // sessionTracker.getCustomProperty = function (key) {
    //     return this[metadataKey].customProperties[key];
    // };

    // sessionTracker.getCustomProperties = function () {
    //     return this[metadataKey].customProperties;
    // };

    // sessionTracker.setCustomProperty = function (key, value) {
    //     this[metadataKey].customProperties[key] = value;

    //     // Create getter and setter for the new custom property
    //     Object.defineProperty(this, key, {
    //         get: function () {
    //             return this[metadataKey].customProperties[key];
    //         },
    //         set: function (newValue) {
    //             this[metadataKey].customProperties[key] = newValue;
    //         }
    //     });
    // };

    Borea.getLastActivityTime = function () {
        const lastActivity = localStorage.getItem('lastActivityTime');
        return lastActivity ? new Date(lastActivity) : null;
    };

    Borea.getSessionData = function () {
        return Object.assign({}, this[metadataKey]);
    };
}

// sessionTracker Maintenance Event Listeners (block for collapsing)
{
    Borea.initMaintenanceEventListeners = function () {
        window.addEventListener('beforeunload', () => {
            this.updateLastActivityTime();
            this.setSessionDuration();
            this.storeMetadataInSessionStorage();
            postData && this.postSessionData();
        });

        window.addEventListener('resize', () => this.updateScreenResolution());

        // TODO not working on testing GitHub page right now
        // TODO multiple getting added
        const eventsForTrackPaths = [
            'popstate',
            'pushstate',
            'replacestate',
        ];
        eventsForTrackPaths.forEach(event => {
            window.addEventListener(event, () => this.captureWindowLocationMetadata());
        });

        // Intercept and handle programmatic navigation
        const originalPushState = history.pushState;
        history.pushState = function () {
            originalPushState.apply(this, arguments);
            Borea.captureWindowLocationMetadata();
        };

        const originalReplaceState = history.replaceState;
        history.replaceState = function () {
            originalReplaceState.apply(this, arguments);
            Borea.captureWindowLocationMetadata();
        };
    };

    Borea.updateLastActivityTime = function () {
        this[metadataKey].lastActivityTime = new Date();
        localStorage.setItem('lastActivityTime', this[metadataKey].lastActivityTime.toISOString());
    };

    Borea.setSessionDuration = function () {
        this[metadataKey].sessionDuration = (new Date() - this[metadataKey].startTime);
    };

    Borea.storeMetadataInSessionStorage = function () {
        sessionStorage.setItem(metadataKey, JSON.stringify(this[metadataKey]));
    };

    // Borea.updateScreenResolution = function () {
    //     this[metadataKey].screenResolution.push(`${window.innerWidth}x${window.innerHeight}`);
    // };

    // Borea.captureWindowLocationMetadata = function () {
    //     const metadata = {
    //         href: window.location.href,
    //         protocol: window.location.protocol,
    //         host: window.location.host,
    //         hostname: window.location.hostname,
    //         port: window.location.port,
    //         pathname: window.location.pathname,
    //         search: window.location.search,
    //         hash: window.location.hash
    //     };

    //     this[metadataKey].userPath.push(metadata.pathname);
    // };

    Borea.postSessionData = function () {
        const route = postSessionDataRoute;

        const isLocalhost = hostname === 'localhost' || hostname === '127.0.0.1';
        const protocol = isLocalhost ? 'http' : 'https';
        const baseUrl = isLocalhost ? `localhost:${GO_PORT}` : DOMAIN;

        const url = `${protocol}://${baseUrl}/${route}`;

        fetch(url, {
            method: 'POST',
            // TODO figure out no-cors client vs server side
            // mode: 'no-cors',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(this.getSessionData()),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data = response.json();
                return data;
            })
            .then(data => {
                console.log('Success:', data);
            })
            .catch(error => {
                console.error('Error:', error);
            });
    };
}

// Event Listener Management
{
    Borea.addEventListener = function ({ target, type, eventListener, options = {}, useCapture, }) {
        if (target === undefined)
            return console.error("ERROR: No target provided to addEventListener to sessionTracker");
        if (type === undefined)
            return console.error("ERROR: No type provided to addEventListener to sessionTracker");

        const controller = new AbortController();
        const eventData = {
            controller,
        };
        options.signal = controller.signal;

        const callback = function (event) {
            if (Borea.events[event.type] === undefined)
                Borea.events[event.type] = [];

            // Push events to arrays
            Borea.events[event.type].push(eventData);
            Borea.eventsArray.push(eventData);

            // Add data to eventData
            eventData.event = event;
            eventData.timeStamp = new Date();

            // console.log(eventData.timeStamp, event.type, event);

            // const defaultEventCallback = sessionTracker.eventTypeMap[type].defaultEventCallback;
            // if (defaultEventCallback != undefined)
            //     defaultEventCallback(event, eventData);

            if (eventListener != undefined)
                eventListener(event, eventData);

            // if (sessionTracker.defaultEventCallback != undefined)
            //     sessionTracker.defaultEventCallback(event, eventData);
            // TODO add custom function to handle eventData
            // if (sessionTracker.options.customEventCallback != undefined)
            //     sessionTracker.options.customEventCallback(event, eventData);
        };

        eventData.eventCallback = callback;

        // Default is false, true makes Event go to listener before EventTarget in DOM tree
        if (options.capture === undefined)
            if (useCapture != undefined)
                options.capture = useCapture;
            else
                options.capture = true;

        try {
            target.addEventListener(type, callback, options);
        } catch (error) {
            console.error(`Failed to add event listener for ${type}:`, error);
        }
    };
}

// init EventListeners on DOM
{
    Borea.initEventListeners = function () {
        // TODO create enabledEventTypes based on options, whitelist and blacklist.
        // if whitelist, then ignore defaultEventTypes and add whitelist first, then blacklist removes from enabledEventTypes.
        // blacklist only, add defaultEventTypes, blacklist used to remove from enabledEventTypes
        // no whitelist no blacklist enabled is same as default

        this.populateEnabledEventTypes();

        this.enabledEventTypes.forEach(type => {
            const params = {};
            params.target = window;
            params.type = type;
            // TODO create eventListener function that is a function that calls all the other functions defined in options. Abstract it to a single function, eventListener

            // TODO pull listener from sessionTracker.options using type as key
            // Also build listener to call universal fn in options and type specific function in options
            // actually add parsing to options that adds them to sessionTracker.eventTypeMap object's by type
            // const customEventListenerOptions = sessionTracker.eventTypeMap[type].customEventListenerOptions;
            // if (customEventListenerOptions != undefined) {
            //     const { eventListener, options, useCapture } = customEventListenerOptions;
            //     if (eventListener != undefined)
            //         params.eventListener = eventListener;
            //     if (options != undefined)
            //         params.options = options;
            //     if (useCapture != undefined)
            //         params.useCapture = useCapture;
            // }

            Borea.addEventListener(params);
        });
    };

    Borea.populateEnabledEventTypes = function () {
        // Object.assign(true, this.enabledEventTypes, this.defaultEventTypes);
        this.enabledEventTypes = this.defaultEventTypes;
    };
}

// tmp data (block for collapsing)
{
    Borea.eventTypeMap = (function () {
        // possibly write get fn's that will check against options and pull default if no options exist
        const eventTypeMap = {
            // Mouse events
            "click": {}, "dblclick": {}, "mousedown": {}, "mouseup": {}, "mousemove": {}, "mouseover": {}, "mouseout": {}, "mouseenter": {}, "mouseleave": {}, "contextmenu": {},
            // Keyboard events
            "keydown": {}, "keyup": {}, "keypress": {},
            // Form events
            "submit": {}, "reset": {}, "change": {}, "input": {}, "invalid": {}, "select": {},
            // Focus events
            "focus": {}, "blur": {}, "focusin": {}, "focusout": {},
            // Window events
            "load": {}, "unload": {}, "beforeunload": {}, "resize": {}, "scroll": {}, "hashchange": {}, "popstate": {},
            // Document and element events
            "DOMContentLoaded": {}, "readystatechange": {}, "cut": {}, "copy": {}, "paste": {},
            // Drag and drop events
            "dragstart": {}, "drag": {}, "dragenter": {}, "dragleave": {}, "dragover": {}, "drop": {}, "dragend": {},
            // Animation and transition events
            "animationstart": {}, "animationend": {}, "animationiteration": {}, "transitionend": {},
            // Media events
            "play": {}, "pause": {}, "ended": {}, "volumechange": {}, "timeupdate": {}, "loadeddata": {}, "canplay": {},
            // Progress events
            "loadstart": {}, "progress": {}, "error": {}, "abort": {}, "load": {}, "loadend": {},
            // Touch events
            "touchstart": {}, "touchmove": {}, "touchend": {}, "touchcancel": {},
            // Pointer events
            "pointerdown": {}, "pointermove": {}, "pointerup": {}, "pointercancel": {}, "pointerover": {}, "pointerout": {}, "pointerenter": {}, "pointerleave": {},
            // Wheel event
            "wheel": {},
            // Storage event
            "storage": {},
            // Server-sent events
            "message": {},
            // Print events
            "beforeprint": {}, "afterprint": {},
            // Clipboard events
            "cut": {}, "copy": {}, "paste": {},
            // Fullscreen events
            "fullscreenchange": {}, "fullscreenerror": {},
            // Visibility events
            "visibilitychange": {},
            // Device orientation and motion events
            "deviceorientation": {}, "devicemotion": {},
            // Page transition events
            "pageshow": {}, "pagehide": {}
        };
        return eventTypeMap;
    })();

    Borea.defaultEventTypes = (function () {
        return [
            // Page lifecycle events
            "load",
            "unload",
            "beforeunload",
            "DOMContentLoaded",
            // User interaction events
            "click",
            "dblclick",
            "contextmenu",
            "mousedown",
            "mouseup",
            // "mousemove",
            "mouseenter",
            "mouseleave",
            "scroll",
            "touchstart",
            "touchend",
            // Form interaction events
            "submit",
            "change",
            "input",
            "focus",
            "blur",
            // Navigation events
            "hashchange",
            "popstate",
            // Visibility and connectivity events
            "visibilitychange",
            "online",
            "offline",
            // Performance and error events
            "loadstart",
            "progress",
            "error",
            "abort",
            "loadend",
            // Media events
            "play",
            "pause",
            "ended"
        ];
    })();
}

// Helpers (block for collapsing)
{
    if (Borea.helpers === undefined)
        Borea.helpers = {};

    Borea.helpers.generateUUID = function () {
        if ('randomUUID' in crypto) {
            return crypto.randomUUID();
        } else {
            // Fallback for older browsers
            return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
                var r = Math.random() * 16 | 0,
                    v = c == 'x' ? r : (r & 0x3 | 0x8);
                return v.toString(16);
            });
        }
    };

    Borea.helpers.getReferrer = function (referrer) {
        if (!referrer) {
            return referrer;
        }

        try {
            const url = new URL(referrer);

            const metadata = {
                protocol: url.protocol,
                hostname: url.hostname,
                port: url.port, // || (url.protocol === 'https:' ? '443' : '80'),
                pathname: url.pathname,
                search: url.search,
                hash: url.hash,
                origin: url.origin,
                host: url.host,
                searchParams: Object.fromEntries(url.searchParams),
                username: url.username,
                password: url.password
            };

            return metadata;
        } catch (error) {
            console.error('Invalid referrer URL:', error);
            return null;
        }
    };

    Borea.helpers.fetchIPAddress = function () {
        return fetch('https://api.ipify.org?format=json',
            // {mode: 'no-cors'}
        )
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log(data.ip)
                return data.ip;
            })
            .catch(error => {
                console.error('Error fetching IP address:', error);
                throw error;
            });
    };
}

Borea.init();
Borea.initEventListeners();

// For CommonJS environments (e.g., Node.js)
if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
    module.exports = Borea;
}

// For ES6 module environments
if (typeof exports !== 'undefined') {
    Object.defineProperty(exports, '__esModule', { value: true });
    exports.default = Borea;
}

} // close file block
