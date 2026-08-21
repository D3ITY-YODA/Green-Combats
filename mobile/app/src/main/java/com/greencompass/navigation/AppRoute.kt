package com.greencompass.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    @Serializable data object Welcome : AppRoute
    @Serializable data object LanguageSelection : AppRoute
    @Serializable data object AccountChoice : AppRoute
    @Serializable data object PersonalRegistration : AppRoute
    @Serializable data object GoogleSignIn : AppRoute
    @Serializable data object PhoneEmailSignIn : AppRoute
    @Serializable data object VerificationCode : AppRoute
    // Future routes will be added here
}
