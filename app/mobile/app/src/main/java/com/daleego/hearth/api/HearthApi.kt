package com.daleego.hearth.api

import retrofit2.http.Body
import retrofit2.http.POST

interface HearthApi {
    @POST("/mobile/pairing")
    suspend fun pairDevice(
        @Body request: PairingRequest
    ): PairingResponse

    @POST("/mobile/locations")
    suspend fun reportLocation(
        @Body request: LocationRequest
    ): LocationResponse
}
