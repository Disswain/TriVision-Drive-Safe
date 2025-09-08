class SosEvent {
  final String carId;
  final double latitude;
  final double longitude;
  final DateTime timestamp;

  SosEvent({required this.carId, required this.latitude, required this.longitude, required this.timestamp});

  factory SosEvent.fromJson(Map<String, dynamic> json) {
    return SosEvent(
      carId: json['car_id'],
      latitude: json['latitude'],
      longitude: json['longitude'],
      timestamp: DateTime.parse(json['timestamp']),
    );
  }
}
