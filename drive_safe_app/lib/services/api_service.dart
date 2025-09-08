import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/parking_session.dart';
import '../models/sos.dart';

class ApiService {
  // For Android emulator use 10.0.2.2; for iOS simulator/macOS/web use localhost.
  static const String baseUrl = 'http://10.0.2.2:8080/api';
  // If you’re not on Android emulator, change to:
  // static const String baseUrl = 'http://localhost:8080/api';

  Future<Map<String, double>> getLocation() async {
    final resp = await http.get(Uri.parse('$baseUrl/location'));
    if (resp.statusCode == 200) {
      final data = jsonDecode(resp.body) as Map<String, dynamic>;
      return {
        'latitude': (data['latitude'] as num).toDouble(),
        'longitude': (data['longitude'] as num).toDouble(),
      };
    }
    throw Exception('Failed to fetch location (${resp.statusCode})');
  }

  Future<ParkingSession> startParking(String carId) async {
    final resp = await http.get(Uri.parse('$baseUrl/parking/start?car_id=$carId'));
    if (resp.statusCode == 200) {
      return ParkingSession.fromJson(jsonDecode(resp.body));
    }
    throw Exception('Failed to start parking (${resp.statusCode})');
  }

  Future<ParkingSession> stopParking(String sessionId) async {
    final resp = await http.get(Uri.parse('$baseUrl/parking/stop?session_id=$sessionId'));
    if (resp.statusCode == 200) {
      return ParkingSession.fromJson(jsonDecode(resp.body));
    }
    throw Exception('Failed to stop parking (${resp.statusCode})');
  }

  Future<List<Map<String, dynamic>>> getNearestParking(double lat, double lng) async {
    final resp = await http.get(Uri.parse('$baseUrl/parking/nearest?lat=$lat&lng=$lng'));
    if (resp.statusCode == 200) {
      return List<Map<String, dynamic>>.from(jsonDecode(resp.body) as List);
    }
    throw Exception('Failed to fetch nearest parking (${resp.statusCode})');
  }

  Future<SosEvent> triggerSos(String carId, double lat, double lng) async {
    final resp = await http.get(Uri.parse('$baseUrl/sos?car_id=$carId&lat=$lat&lng=$lng'));
    if (resp.statusCode == 200) {
      return SosEvent.fromJson(jsonDecode(resp.body));
    }
    throw Exception('Failed to send SOS (${resp.statusCode})');
  }
}
