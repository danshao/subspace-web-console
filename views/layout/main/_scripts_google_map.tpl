<script>
  function initMap() {
    var locations = [ 
      {{range.location}} 
      { 
        lat: {{.Lat}},
        lng: {{.Lng}},
        info: {{ .Info }} 
      },
      {{ end }}
    ]

    var mapOptions = {
      zoom: 2,
      streetViewControl : false,
      mapTypeControl: false,
      rotateControl: false,
      fullScreenControl: false,
      minZoom: 2,
      maxZoom: 9,
      zoom: 2,
      center: {
        lat: 10,
        lng: 0
      },
      styles:
      [{"featureType":"administrative.country","elementType":"geometry","stylers":[{"color":"#191c2d"},{"weight":1.5}]},{"featureType":"administrative.country","elementType":"labels","stylers":[{"color":"#191c2d"},{"lightness":"58"}]},{"featureType":"administrative.country","elementType":"labels.text.fill","stylers":[{"lightness":80}]},{"featureType":"administrative.country","elementType":"labels.text.stroke","stylers":[{"color":"#252938"}]},{"featureType":"administrative.land_parcel","stylers":[{"visibility":"off"}]},{"featureType":"administrative.locality","stylers":[{"visibility":"on"}]},{"featureType":"administrative.locality","elementType":"labels.text","stylers":[{"lightness":55},{"visibility":"on"}]},{"featureType":"administrative.locality","elementType":"labels.text.fill","stylers":[{"color":"#d3d8e6"},{"lightness":-20}]},{"featureType":"administrative.locality","elementType":"labels.text.stroke","stylers":[{"color":"#191d2d"},{"visibility":"on"}]},{"featureType":"administrative.neighborhood","stylers":[{"visibility":"off"}]},{"featureType":"administrative.province","stylers":[{"visibility":"on"}]},{"featureType":"administrative.province","elementType":"labels.text","stylers":[{"visibility":"on"}]},{"featureType":"administrative.province","elementType":"labels.text.fill","stylers":[{"color":"#8792ba"},{"lightness":30}]},{"featureType":"administrative.province","elementType":"labels.text.stroke","stylers":[{"color":"#252938"},{"lightness":-65}]},{"featureType":"landscape","stylers":[{"color":"#39415e"},{"visibility":"simplified"}]},{"featureType":"poi","stylers":[{"visibility":"off"}]},{"featureType":"road","stylers":[{"color":"#28314f"},{"lightness":"4"},{"visibility":"off"}]},{"featureType":"road","elementType":"labels","stylers":[{"visibility":"off"}]},{"featureType":"transit","stylers":[{"visibility":"off"}]},{"featureType":"water","stylers":[{"color":"#191c2d"},{"visibility":"on"}]}]
    }
    var map = new google.maps.Map(document.getElementById('map'), mapOptions);
    var bounds = new google.maps.LatLngBounds();
    var locationDic = {};

    //TODO: Ideal location data structure - not currently implemented
    // {
    //   "1, 1": {
    //     0: {
    //       info: "dan@me.com",
    //       count: 1
    //     }
    //     1: {
    //       info: "123@123.com",
    //       count: 2
    //     }
    //   }
    // }

    // Merge locations if lat/lng are the same
    for (var i = 0; i < locations.length; i++) {
      var key = String(locations[i].lat) + ',' + String(locations[i].lng)
      if (locationDic[key] === undefined) {
        locationDic[key] = locations[i].info
      } else {
        if (locationDic[key].search(locations[i].info) == -1) {
          locationDic[key] = locationDic[key] + '<br>' + locations[i].info
        }
      }
    }

    $.each(locationDic, function (index, value) {
      lat = index.split(",")[0]
      lng = index.split(",")[1]

      // Initialise the infoWindow
      infoWindow = new google.maps.InfoWindow({});
      
      // Add a marker to the map based on our coordinates
      if (value == "server") {
        var marker = new google.maps.Marker({
          position: new google.maps.LatLng(lat, lng),
          map: map,
          icon: customIcon({
            fillColor: '#09ee05',
            strokeColor: '#000'
          }),
          html: "Subspace Server"
        })
      } else {
        var marker = new google.maps.Marker({
            position: new google.maps.LatLng(lat, lng),
            map: map,
            html: value
        });
      }

      // Display our info window when the marker is clicked
      google.maps.event.addListener(marker, 'click', function() {
        infoWindow.setContent(this.html);
        infoWindow.open(map, this);
      });
      // extend map bounds
      bounds.extend(marker.getPosition());

      // fit map to bounds
      map.fitBounds(bounds);
    });
  }

  function customIcon (opts) {
    return Object.assign({
      path: 'M 0,0 C -2,-20 -10,-22 -10,-30 A 10,10 0 1,1 10,-30 C 10,-22 2,-20 0,0 z M -2,-30 a 2,2 0 1,1 4,0 2,2 0 1,1 -4,0',
      fillColor: '#34495e',
      fillOpacity: 1,
      strokeColor: '#000',
      strokeWeight: 1,
      scale: 1,
    }, opts);
  }
</script>
<script async defer src="https://maps.googleapis.com/maps/api/js?key=AIzaSyCrfTkjmt9Nq0WvcxelBMt-SUdIAKniHM0&callback=initMap"></script>