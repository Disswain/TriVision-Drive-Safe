class ParkingSession {
  final String sessionId;
  final String carId;
  final DateTime startedAt;
  final DateTime? stoppedAt;
  final int? durationMinutes;

  ParkingSession({
    required this.sessionId,
    required this.carId,
    required this.startedAt,
    this.stoppedAt,
    this.durationMinutes,
  });

  factory ParkingSession.fromJson(Map<String, dynamic> json) {
    return ParkingSession(
      sessionId: json['session_id'],
      carId: json['car_id'],
      startedAt: DateTime.parse(json['started_at']),
      stoppedAt:
          json['stopped_at'] != null ? DateTime.parse(json['stopped_at']) : null,
      durationMinutes: json['duration_minutes'],
    );
  }
}
