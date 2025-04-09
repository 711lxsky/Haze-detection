interface GeolocationCoordinates {
  latitude: number;
  longitude: number;
  accuracy: number;
  timestamp: number;
}

interface GeolocationOptions {
  enableHighAccuracy?: boolean;
  timeout?: number;
  maximumAge?: number;
}

function getGeolocation(
  successCallback: (coords: GeolocationCoordinates) => void,
  errorCallback?: (error: string) => void,
  options: GeolocationOptions = {}
): void {
  // 检查浏览器是否支持地理定位
  if (!navigator.geolocation) {
    if (errorCallback) {
      errorCallback('您的浏览器不支持地理定位功能');
    }
    return;
  }

  // 默认选项
  const defaultOptions: GeolocationOptions = {
    enableHighAccuracy: true,
    timeout: 10000,
    maximumAge: 0
  };

  // 合并选项
  const geoOptions: GeolocationOptions = { ...defaultOptions, ...options };

  // 成功获取位置的回调
  const onSuccess = (position: GeolocationPosition): void => {
    const coords: GeolocationCoordinates = {
      latitude: position.coords.latitude,
      longitude: position.coords.longitude,
      accuracy: position.coords.accuracy,
      timestamp: position.timestamp
    };

    successCallback(coords);
  };

  // 获取位置失败的回调
  const onError = (error: GeolocationPositionError): void => {
    let errorMessage: string = '';

    switch (error.code) {
      case error.PERMISSION_DENIED:
        errorMessage = '用户拒绝了地理定位请求';
        break;
      case error.POSITION_UNAVAILABLE:
        errorMessage = '位置信息不可用';
        break;
      case error.TIMEOUT:
        errorMessage = '获取用户位置超时';
        break;
      default:
        errorMessage = '发生未知错误';
        break;
    }

    if (errorCallback) {
      errorCallback(errorMessage);
    }
  };

  // 获取位置
  navigator.geolocation.getCurrentPosition(
    onSuccess,
    onError,
    geoOptions as PositionOptions
  );
}

function getGeolocationPromise(
  options: GeolocationOptions = {}
): Promise<GeolocationCoordinates> {
  return new Promise((resolve, reject) => {
    getGeolocation(
      (coords) => resolve(coords),
      (error) => reject(new Error(error)),
      options
    );
  });
}

export {
  getGeolocationPromise,
}