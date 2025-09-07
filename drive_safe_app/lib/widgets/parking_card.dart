import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class ParkingCard extends StatelessWidget {
  final Map<String, dynamic> spot;
  const ParkingCard({super.key, required this.spot});

  Future<void> _openMaps() async {
    final lat = (spot['lat'] as num).toDouble();
    final lng = (spot['lng'] as num).toDouble();
    final url = Uri.parse(
      'https://www.google.com/maps/dir/?api=1&destination=$lat,$lng',
    );
    await launchUrl(url, mode: LaunchMode.externalApplication);
  }

  @override
  Widget build(BuildContext context) {
    final distance = (spot['distance_km'] as num?)?.toDouble();
    final slots = spot['available_slots'];
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: ListTile(
        title: Text(spot['name']?.toString() ?? 'Parking'),
        subtitle: Text(
          'Slots: $slots • ${distance != null ? '${distance.toStringAsFixed(2)} km' : 'distance n/a'}',
        ),
        trailing: IconButton(
          icon: const Icon(Icons.directions),
          onPressed: _openMaps,
        ),
      ),
    );
  }
}
