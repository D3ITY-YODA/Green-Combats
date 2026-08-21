package com.greencompass.core.ui

import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.CloudOff
import androidx.compose.material.icons.outlined.ErrorOutline
import androidx.compose.material.icons.outlined.Info
import androidx.compose.material.icons.outlined.LocationOff
import androidx.compose.material.icons.outlined.Lock
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp

@Composable
fun LoadingSkeleton(modifier: Modifier = Modifier) {
    val infiniteTransition = rememberInfiniteTransition()
    val alpha by infiniteTransition.animateFloat(
        initialValue = 0.3f, targetValue = 0.7f,
        animationSpec = infiniteRepeatable(animation = tween(1000, easing = LinearEasing))
    )
    Box(
        modifier = modifier.clip(RoundedCornerShape(12.dp)).background(GreenCompassColors.Stone.copy(alpha = alpha))
    )
}

@Composable
fun EmptyStateView(icon: ImageVector, title: String, message: String, modifier: Modifier = Modifier) {
    Column(modifier = modifier.fillMaxWidth().padding(AppSpacing.xxl), horizontalAlignment = Alignment.CenterHorizontally) {
        Icon(imageVector = icon, contentDescription = null, tint = GreenCompassColors.Sage, modifier = Modifier.size(64.dp))
        Spacer(modifier = Modifier.height(AppSpacing.md))
        Text(text = title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xs))
        Text(text = message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
    }
}

@Composable
fun ErrorStateView(message: String, onRetry: () -> Unit, modifier: Modifier = Modifier) {
    Column(modifier = modifier.fillMaxWidth().padding(AppSpacing.xxl), horizontalAlignment = Alignment.CenterHorizontally) {
        Icon(imageVector = Icons.Outlined.ErrorOutline, contentDescription = null, tint = GreenCompassColors.EmergencyRed, modifier = Modifier.size(64.dp))
        Spacer(modifier = Modifier.height(AppSpacing.md))
        Text(text = "Something went wrong", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xs))
        Text(text = message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xl))
        PrimaryButton(text = "Try again", onClick = onRetry, modifier = Modifier.width(200.dp))
    }
}

@Composable
fun OfflineBanner(modifier: Modifier = Modifier) {
    Row(
        modifier = modifier.fillMaxWidth().background(GreenCompassColors.Mist).padding(AppSpacing.md),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(imageVector = Icons.Outlined.CloudOff, contentDescription = null, tint = GreenCompassColors.MutedText, modifier = Modifier.size(20.dp))
        Spacer(modifier = Modifier.width(AppSpacing.sm))
        Text(text = "You're offline. Showing information saved on this device.", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
    }
}

@Composable
fun PermissionDeniedStateView(
    title: String = "Location permission is off",
    message: String = "You can still choose a place manually.",
    onChoosePlace: () -> Unit,
    onOpenSettings: () -> Unit
) {
    Column(modifier = Modifier.fillMaxWidth().padding(AppSpacing.xxl), horizontalAlignment = Alignment.CenterHorizontally) {
        Icon(imageVector = Icons.Outlined.LocationOff, contentDescription = null, tint = GreenCompassColors.MutedText, modifier = Modifier.size(64.dp))
        Spacer(modifier = Modifier.height(AppSpacing.md))
        Text(text = title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xs))
        Text(text = message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xl))
        PrimaryButton(text = "Choose a place", onClick = onChoosePlace, modifier = Modifier.padding(bottom = AppSpacing.sm))
        TextLinkButton(text = "Open settings", onClick = onOpenSettings)
    }
}

@Composable
fun SessionExpiredStateView(onSignIn: () -> Unit) {
    Column(modifier = Modifier.fillMaxWidth().padding(AppSpacing.xxl), horizontalAlignment = Alignment.CenterHorizontally) {
        Icon(imageVector = Icons.Outlined.Lock, contentDescription = null, tint = GreenCompassColors.MutedText, modifier = Modifier.size(64.dp))
        Spacer(modifier = Modifier.height(AppSpacing.md))
        Text(text = "Your session has ended", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xs))
        Text(text = "Please sign in again to continue.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
        Spacer(modifier = Modifier.height(AppSpacing.xl))
        PrimaryButton(text = "Sign in", onClick = onSignIn, modifier = Modifier.width(200.dp))
    }
}
