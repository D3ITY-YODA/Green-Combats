package com.greencompass.data.mapper

import com.greencompass.data.local.entity.UpdateEntity
import com.greencompass.data.remote.dto.UpdateDto

fun UpdateDto.toEntity(): UpdateEntity {
    return UpdateEntity(
        id = this.id,
        placeId = this.placeId,
        placeName = this.placeName,
        topicKey = this.topicKey,
        type = this.type,
        priority = this.priority,
        title = this.title,
        message = this.message,
        validFrom = this.validFrom,
        validUntil = this.validUntil,
        updatedAt = this.updatedAt,
        sourceName = this.sourceName
    )
}
