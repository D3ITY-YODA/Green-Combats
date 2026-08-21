package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RequestAccessScreen(
    onBack: () -> Unit,
    onRequestSent: () -> Unit
) {
    var reason by remember { mutableStateOf("") }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Request access",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            OutlinedTextField(
                value = "Lower Valley Water Authority",
                onValueChange = {},
                label = { Text("Organization") },
                readOnly = true,
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.md),
                shape = RoundedCornerShape(12.dp)
            )

            OutlinedTextField(
                value = reason,
                onValueChange = { reason = it },
                label = { Text("Reason for joining") },
                placeholder = { Text("Optional") },
                modifier = Modifier.fillMaxWidth().height(120.dp).padding(bottom = AppSpacing.xl),
                shape = RoundedCornerShape(12.dp)
            )

            Spacer(modifier = Modifier.weight(1f))

            PrimaryButton(
                text = "Send request",
                onClick = onRequestSent
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
