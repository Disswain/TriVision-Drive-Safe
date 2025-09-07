import 'package:flutter/material.dart';

class LocationTile extends StatelessWidget {
  final double latitude;
  final double longitude;
  const LocationTile({super.key, required this.latitude, required this.longitude});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: const Icon(Icons.location_on),
      title: Text('Lat: ${latitude.toStringAsFixed(6)}'),
      subtitle: Text('Lng: ${longitude.toStringAsFixed(6)}'),
    );
  }
}
