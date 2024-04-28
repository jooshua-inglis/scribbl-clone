export interface BaseEvent {
  type: string | number;
}

export interface Event extends BaseEvent {
  [attachment: string]: unknown;
}

export type EventHandler<E, T> = (event: E & { type: T }) => void;

export class PubSub<T extends BaseEvent = Event> {
  private _listeners: { [key: string]: Array<CallableFunction> };

  constructor() {
    this._listeners = {};
  }

  addEventListener<E extends T["type"]>(
    type: E,
    listener: EventHandler<T, E>
  ) {
    if (this._listeners === undefined) this._listeners = {};

    const listeners = this._listeners;

    if (listeners[type] === undefined) {
      listeners[type] = [];
    }

    if (listeners[type].indexOf(listener) === -1) {
      listeners[type].push(listener);
    }
    return () => {
      this.removeEventListener(type, listener)
    }
  }

  hasEventListener<E extends T["type"]>(
    type: E,
    listener: EventHandler<T, E>
  ) {
    if (this._listeners === undefined) return false;

    const listeners = this._listeners;

    return (
      listeners[type] !== undefined && listeners[type].indexOf(listener) !== -1
    );
  }

  removeEventListener<E extends T["type"]>(
    type: E,
    listener: EventHandler<T, E>
  ) {
    if (this._listeners === undefined) return;

    const listeners = this._listeners;
    const listenerArray = listeners[type];

    if (listenerArray !== undefined) {
      const index = listenerArray.indexOf(listener);

      if (index !== -1) {
        listenerArray.splice(index, 1);
      }
    }
  }

  dispatchEvent(event: T) {
      if (this._listeners === undefined) return;
      const listeners = this._listeners;
      const listenerArray = listeners[event.type];

    if (listenerArray !== undefined) {
      // Make a copy, in case listeners are removed while iterating.
      const array = listenerArray.slice(0);

      for (let i = 0, l = array.length; i < l; i++) {
        (async () => array[i](event))();
      }
    }
  }
}
