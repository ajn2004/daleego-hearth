package com.daleego.hearth.api

import com.squareup.moshi.Json

data class PairingRequest(
    @Json(name = "pair_code")
    val pairCode: String,

    val platform: String = "android",

    val model: String,

    @Json(name = "model_type")
    val modelType: String
)

data class PairingResponse(
    @Json(name = "api_key")
    val apiKey: String,

    val device: DeviceDto
)

data class DeviceDto(
    val id: String,

    @Json(name = "person_id")
    val personId: String?,

    val name: String,
    val platform: String
)

data class LocationRequest(
    val lat: Double,
    val lng: Double,

    @Json(name = "accuracy")
    val accuracy: Double,

    @Json(name = "recorded_at")
    val recordedAt: String
)

data class LocationResponse(
    val id: String? = null
)