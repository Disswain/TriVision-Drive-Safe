class CarLocation {
  final double latitude;
  final double longitude;

  CarLocation({required this.latitude, required this.longitude});

  factory CarLocation.fromJson(Map<String, dynamic> json) {
    return CarLocation(
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
    );
  }
}
